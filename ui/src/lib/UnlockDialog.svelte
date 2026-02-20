<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { UnlockVault } from '../../wailsjs/go/main/App.js';

    const dispatch = createEventDispatcher();

    let password = '';
    let error = '';
    let loading = false;
    let attempts = 0;
    let showPassword = false;

    async function handleSubmit() {
        error = '';

        if (!password) {
            error = 'Please enter your master password.';
            return;
        }

        loading = true;

        try {
            await UnlockVault(password);
            dispatch('unlock');
        } catch (e: any) {
            attempts++;
            password = '';
            error = 'Incorrect master password.';

            if (attempts >= 3) {
                error = 'Incorrect master password. Please try again carefully.';
            }
        } finally {
            loading = false;
        }
    }

    function handleKeydown(e: KeyboardEvent) {
        if (e.key === 'Enter' && password && !loading) {
            handleSubmit();
        }
    }
</script>

<div class="vault-dialog" on:keydown={handleKeydown}>
    <div class="vault-icon">
        <svg viewBox="0 0 24 24" width="48" height="48" fill="none" stroke="currentColor" stroke-width="1.5">
            <rect x="3" y="11" width="18" height="11" rx="2" ry="2" />
            <path d="M7 11V7a5 5 0 0 1 10 0v4" />
            <circle cx="12" cy="16.5" r="1.5" />
        </svg>
    </div>

    <h2>Unlock GitSyncer</h2>
    <p class="subtitle">Enter your master password to access credentials.</p>

    <div class="field">
        <div class="input-wrap">
            {#if showPassword}
                <input
                    type="text"
                    placeholder="Master password"
                    bind:value={password}
                    class:field-error={error !== ''}
                    autocomplete="current-password"
                />
            {:else}
                <input
                    type="password"
                    placeholder="Master password"
                    bind:value={password}
                    class:field-error={error !== ''}
                    autocomplete="current-password"
                />
            {/if}
            <button
                class="toggle-vis"
                type="button"
                tabindex="-1"
                on:click={() => (showPassword = !showPassword)}
                title={showPassword ? 'Hide password' : 'Show password'}
            >
                {#if showPassword}
                    <svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94" />
                        <path d="M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19" />
                        <line x1="1" y1="1" x2="23" y2="23" />
                    </svg>
                {:else}
                    <svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" />
                        <circle cx="12" cy="12" r="3" />
                    </svg>
                {/if}
            </button>
        </div>
    </div>

    {#if error}
        <div class="error-banner">{error}</div>
    {/if}

    <button class="btn-primary" on:click={handleSubmit} disabled={!password || loading}>
        {#if loading}Unlocking...{:else}Unlock{/if}
    </button>
</div>
