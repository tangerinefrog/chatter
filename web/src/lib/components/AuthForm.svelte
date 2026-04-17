<script lang="ts">
    import '$lib/css/auth.css';

    export let error = '';
    export let isLoading = false;
    export let onSubmit: (request: { username: string; password: string; mode: 'login' | 'signup' }) => void;
    export let onToggleMode: () => void;

    let mode: 'login' | 'signup' = 'login';
    let username = '';
    let password = '';

    function handleSubmit() {
        const trimmedUsername = username.trim();
        const trimmedPassword = password.trim();

        if (!trimmedUsername || !trimmedPassword) {
            return;
        }

        onSubmit?.({ username: trimmedUsername, password: trimmedPassword, mode });
    }

    function handleKeydown(event: KeyboardEvent) {
        if (event.key === 'Enter') {
            event.preventDefault();
            handleSubmit();
        }
    }

    function toggleMode() {
        mode = mode === 'login' ? 'signup' : 'login';
        username = '';
        password = '';
        onToggleMode?.();
    }
</script>

<style>
    .container {
        width: 100%;
    }
</style>

<div class="container">
    <div class="card">
        <h2>{mode === 'login' ? 'Login' : 'Create Account'}</h2>

        <form on:submit|preventDefault={handleSubmit}>
            <input
                class="input"
                type="text"
                placeholder="Username"
                bind:value={username}
                on:keydown={handleKeydown}
                aria-label="Username"
                required
                disabled={isLoading}
                autocomplete="username"
            />

            <input
                class="input"
                type="password"
                placeholder="Password"
                bind:value={password}
                on:keydown={handleKeydown}
                aria-label="Password"
                required
                disabled={isLoading}
                autocomplete={mode === 'login' ? 'current-password' : 'new-password'}
            />

            <button
                class="button"
                type="submit"
                disabled={isLoading || !username.trim() || !password.trim()}
            >
                {#if isLoading}
                    {mode === 'login' ? 'Logging in...' : 'Signing up...'}
                {:else}
                    {mode === 'login' ? 'Login' : 'Sign Up'}
                {/if}
            </button>

            {#if error}
                <div class="error" role="alert">{error}</div>
            {/if}
        </form>

        <p class="switch">
            {mode === 'login'
                ? "Don't have an account?"
                : "Already have an account?"}
            <button
                type="button"
                on:click={toggleMode}
                disabled={isLoading}
            >
                {mode === 'login' ? 'Sign up' : 'Login'}
            </button>
        </p>
    </div>
</div>
