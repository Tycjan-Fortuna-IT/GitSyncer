<script lang="ts">
    import { onDestroy, onMount } from 'svelte';
    import {
        IsMasterPasswordSetup,
        LockVault,
    } from '../wailsjs/go/main/App.js';
    import MasterPasswordSetup from './lib/MasterPasswordSetup.svelte';
    import UnlockDialog from './lib/UnlockDialog.svelte';

    enum AppState {
        Loading = 'loading',
        Setup = 'setup',
        Locked = 'locked',
        Unlocked = 'unlocked',
    }

    let state: AppState = AppState.Loading;
    let lockTimer: ReturnType<typeof setTimeout> | null = null;

    const LOCK_TIMEOUT_MS = 15 * 60 * 1000; // 15 minutes

    onMount(async () => {
        try {
            const isSetup = await IsMasterPasswordSetup();
            state = isSetup ? AppState.Locked : AppState.Setup;
        } catch {
            state = AppState.Setup;
        }
    });

    onDestroy(() => {
        clearLockTimer();
    });

    function onSetupComplete() {
        state = AppState.Unlocked;
        startLockTimer();
    }

    function onUnlocked() {
        state = AppState.Unlocked;
        startLockTimer();
    }

    async function lockVault() {
        clearLockTimer();

        try {
            await LockVault();
        } catch {
            // Lock locally regardless
        }

        state = AppState.Locked;
    }


    function startLockTimer() {
        clearLockTimer();
        lockTimer = setTimeout(() => lockVault(), LOCK_TIMEOUT_MS);
    }

    function clearLockTimer() {
        if (lockTimer !== null) {
            clearTimeout(lockTimer);
            lockTimer = null;
        }
    }

    function resetLockTimer() {
        if (state === AppState.Unlocked) {
            startLockTimer();
        }
    }
</script>

<svelte:window
    on:mousemove={resetLockTimer}
    on:mousedown={resetLockTimer}
    on:keydown={resetLockTimer}
    on:scroll={resetLockTimer}
    on:touchstart={resetLockTimer}
/>

<main>
    {#if state === AppState.Loading}
        <div class="vault-dialog">
            <div class="loading-spinner"></div>
            <p>Loading...</p>
        </div>
    {:else if state === AppState.Setup}
        <MasterPasswordSetup on:setup={onSetupComplete} />
    {:else if state === AppState.Locked}
        <UnlockDialog on:unlock={onUnlocked} />
    {:else if state === AppState.Unlocked}
        <div class="app-content">
            <header class="app-header">
                <h1>GitSyncer</h1>
                <button class="btn-lock" on:click={lockVault} title="Lock vault">
                    <svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2">
                        <rect x="3" y="11" width="18" height="11" rx="2" ry="2" />
                        <path d="M7 11V7a5 5 0 0 1 10 0v4" />
                    </svg>
                    Lock
                </button>
            </header>

            <div class="main-area">
                <p class="status-text">Vault unlocked. Credentials are accessible.</p>
            </div>
        </div>
    {/if}
</main>

<style>
    main {
        display: flex;
        align-items: center;
        justify-content: center;
        min-height: 100vh;
        padding: 2rem;
        box-sizing: border-box;
    }

    :global(.vault-dialog) {
        background: rgba(255, 255, 255, 0.05);
        border: 1px solid rgba(255, 255, 255, 0.1);
        border-radius: 12px;
        padding: 2.5rem 2rem;
        max-width: 380px;
        width: 100%;
        text-align: center;
    }

    :global(.vault-icon) {
        color: rgba(255, 255, 255, 0.6);
        margin-bottom: 1rem;
    }

    :global(.vault-dialog h2) {
        margin: 0 0 0.5rem;
        font-size: 1.4rem;
        font-weight: 600;
    }

    :global(.subtitle) {
        color: rgba(255, 255, 255, 0.5);
        font-size: 0.85rem;
        margin: 0 0 1.5rem;
        line-height: 1.4;
    }

    :global(.field) {
        margin-bottom: 1rem;
        text-align: left;
    }

    :global(.input-wrap) {
        position: relative;
    }

    :global(.field input::-ms-reveal) {
        display: none;
    }

    :global(.field input) {
        width: 100%;
        padding: 0.6rem 2.5rem 0.6rem 0.75rem;
        border: 1px solid rgba(255, 255, 255, 0.15);
        border-radius: 6px;
        background: rgba(0, 0, 0, 0.25);
        color: white;
        font-size: 0.95rem;
        font-family: inherit;
        outline: none;
        box-sizing: border-box;
        transition: border-color 0.15s;
    }

    :global(.field input:focus) {
        border-color: rgba(100, 160, 255, 0.5);
    }

    :global(.toggle-vis) {
        position: absolute;
        right: 0.5rem;
        top: 50%;
        transform: translateY(-50%);
        background: none;
        border: none;
        padding: 0.25rem;
        cursor: pointer;
        color: rgba(255, 255, 255, 0.35);
        display: flex;
        align-items: center;
        transition: color 0.15s;
    }

    :global(.toggle-vis:hover) {
        color: rgba(255, 255, 255, 0.7);
    }

    :global(.field input.field-error) {
        border-color: rgba(255, 100, 100, 0.6);
    }

    :global(.field .hint) {
        display: block;
        font-size: 0.75rem;
        margin-top: 0.3rem;
        color: rgba(255, 255, 255, 0.4);
    }

    :global(.error-text) {
        color: rgba(255, 120, 120, 0.9) !important;
    }

    :global(.error-banner) {
        background: rgba(255, 80, 80, 0.15);
        border: 1px solid rgba(255, 80, 80, 0.3);
        border-radius: 6px;
        padding: 0.5rem 0.75rem;
        margin-bottom: 1rem;
        font-size: 0.85rem;
        color: rgba(255, 160, 160, 1);
    }

    :global(.btn-primary) {
        width: 100%;
        padding: 0.65rem;
        border: none;
        border-radius: 6px;
        background: rgba(60, 130, 240, 0.85);
        color: white;
        font-size: 0.95rem;
        font-family: inherit;
        cursor: pointer;
        transition: background 0.15s;
    }

    :global(.btn-primary:hover:not(:disabled)) {
        background: rgba(80, 145, 255, 0.95);
    }

    :global(.btn-primary:disabled) {
        opacity: 0.4;
        cursor: not-allowed;
    }

    .loading-spinner {
        width: 32px;
        height: 32px;
        margin: 0 auto 1rem;
        border: 3px solid rgba(255, 255, 255, 0.15);
        border-top-color: rgba(255, 255, 255, 0.6);
        border-radius: 50%;
        animation: spin 0.7s linear infinite;
    }

    @keyframes spin {
        to { transform: rotate(360deg); }
    }

    .app-content {
        width: 100%;
        max-width: 900px;
        align-self: flex-start;
    }

    .app-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding-bottom: 1rem;
        border-bottom: 1px solid rgba(255, 255, 255, 0.1);
        margin-bottom: 2rem;
    }

    .app-header h1 {
        margin: 0;
        font-size: 1.5rem;
        font-weight: 600;
    }

    .btn-lock {
        display: flex;
        align-items: center;
        gap: 0.4rem;
        padding: 0.4rem 0.75rem;
        border: 1px solid rgba(255, 255, 255, 0.15);
        border-radius: 6px;
        background: transparent;
        color: rgba(255, 255, 255, 0.7);
        font-size: 0.85rem;
        font-family: inherit;
        cursor: pointer;
        transition: border-color 0.15s, color 0.15s;
    }

    .btn-lock:hover {
        border-color: rgba(255, 255, 255, 0.3);
        color: white;
    }

    .main-area {
        text-align: center;
    }

    .status-text {
        color: rgba(255, 255, 255, 0.5);
        font-size: 0.95rem;
    }
</style>
