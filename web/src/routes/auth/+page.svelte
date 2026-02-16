<script lang="ts">
    import { goto } from '$app/navigation';
    import { login, signup, type ApiError } from '$lib/api/client';
    import '$lib/css/auth.css';

    let mode: 'login' | 'signup' = 'login';
    let username = '';
    let password = '';
    let error = '';
    let isLoading = false;
    
    async function submit() {
        const trimmedUsername = username.trim();
        const trimmedPassword = password.trim();

        if (!trimmedUsername) {
            error = 'Username is required';
            return;
        }

        if (!trimmedPassword) {
            error = 'Password is required';
            return;
        }

        error = '';
        isLoading = true;

        try {
            if (mode === 'login') {
                await login({ username: trimmedUsername, password: trimmedPassword });
            } else {
                await signup({ username: trimmedUsername, password: trimmedPassword });
            }
            goto('/');
        } catch (err) {
            const apiError = err as ApiError;
            error = apiError.message || `Failed to ${mode}`;
        } finally {
            isLoading = false;
        }
    }

    function toggleMode() {
        mode = mode === 'login' ? 'signup' : 'login';
        error = '';
        username = '';
        password = '';
    }

    function handleKeydown(event: KeyboardEvent) {
        if (event.key === 'Enter') {
            event.preventDefault();
            submit();
        }
    }
</script>

<div class="container">
    <div class="card">
        <h2>{mode === 'login' ? 'Login' : 'Create Account'}</h2>

        <form on:submit|preventDefault={submit}>
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
                on:click={toggleMode}
                disabled={isLoading}
            >
                {mode === 'login' ? 'Sign up' : 'Login'}
            </button>
        </p>
    </div>
</div>