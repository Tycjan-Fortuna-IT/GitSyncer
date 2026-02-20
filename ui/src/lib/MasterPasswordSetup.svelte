<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { SetupMasterPassword } from '../../wailsjs/go/main/App.js';

    const dispatch = createEventDispatcher();

    let password = '';
    let confirmPassword = '';
    let error = '';
    let loading = false;
    let showPassword = false;
    let showConfirm = false;

    $: mismatch = confirmPassword.length > 0 && password !== confirmPassword;
    $: tooShort = password.length > 0 && password.length < 8;
    $: canSubmit = password.length >= 8 && password === confirmPassword && !loading;

    async function handleSubmit() {
        error = '';

        if (password !== confirmPassword) {
            error = 'Passwords do not match.';
            return;
        }

        if (password.length < 8) {
            error = 'Password must be at least 8 characters.';
            return;
        }

        loading = true;

        try {
            await SetupMasterPassword(password);
            dispatch('setup');
        } catch (e: any) {
            error = e?.message || String(e) || 'Failed to set master password.';
        } finally {
            loading = false;
        }
    }

    function handleKeydown(e: KeyboardEvent) {
        if (e.key === 'Enter' && canSubmit) {
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

    <h2>Create Master Password</h2>
    <p class="subtitle">
        This password encrypts all your stored credentials.
        It cannot be recovered if lost.
    </p>

    <div class="field">
        <div class="input-wrap">
            {#if showPassword}
                <input
                    type="text"
                    placeholder="Master password"
                    bind:value={password}
                    class:field-error={tooShort}
                    autocomplete="new-password"
                />
            {:else}
                <input
                    type="password"
                    placeholder="Master password"
                    bind:value={password}
                    class:field-error={tooShort}
                    autocomplete="new-password"
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
        {#if tooShort}
            <span class="hint error-text">Must be at least 8 characters</span>
        {/if}
    </div>

    <div class="field">
        <div class="input-wrap">
            {#if showConfirm}
                <input
                    type="text"
                    placeholder="Confirm password"
                    bind:value={confirmPassword}
                    class:field-error={mismatch}
                    autocomplete="new-password"
                />
            {:else}
                <input
                    type="password"
                    placeholder="Confirm password"
                    bind:value={confirmPassword}
                    class:field-error={mismatch}
                    autocomplete="new-password"
                />
            {/if}
            <button
                class="toggle-vis"
                type="button"
                tabindex="-1"
                on:click={() => (showConfirm = !showConfirm)}
                title={showConfirm ? 'Hide password' : 'Show password'}
            >
                {#if showConfirm}
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
        {#if mismatch}
            <span class="hint error-text">Passwords do not match</span>
        {/if}
    </div>

    {#if error}
        <div class="error-banner">{error}</div>
    {/if}

    <button class="btn-primary" on:click={handleSubmit} disabled={!canSubmit}>
        {#if loading}Setting up...{:else}Set Master Password{/if}
    </button>
</div>
