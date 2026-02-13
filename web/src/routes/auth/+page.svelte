<script lang="ts">
    import { goto } from '$app/navigation';
    import '$lib/css/auth.css';

    let mode: 'login' | 'signup' = 'login';
    let username = '';
    let password = '';
    let error = '';

    const API_URL = 'http://localhost:8080/api/auth';

    async function submit() {
        error = '';

        try {
            const res = await fetch(`${API_URL}/${mode}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ username, password }),
                credentials: 'include'
            });

            if (!res.ok) {
                const data = await res.json();
                throw new Error(data.message || 'Error occured');
            }

            goto('/');
        } catch (err: any) {
            error = err.message;
        }
    }

    function toggleMode() {
        mode = mode === 'login' ? 'signup' : 'login';
        error = '';
    }
</script>

<div class="container">
    <div class="card">
        <h2>{mode === 'login' ? 'Login' : 'Create Account'}</h2>

        <input
            class="input"
            type="text"
            placeholder="Username"
            bind:value={username}
        />

        <input
            class="input"
            type="password"
            placeholder="Password"
            bind:value={password}
        />

        {#if error}
            <div class="error">{error}</div>
        {/if}

        <button 
            class="button"
            on:click={submit}>
            {mode === 'login' ? 'Login' : 'Sign Up'}
        </button>

        <p class="switch">
            {mode === 'login'
                ? "Don't have an account?"
                : "Already have an account?"}
            <button 
                on:click={toggleMode}>
                {mode === 'login' ? 'Sign up' : 'Login'}
            </button>
        </p>
    </div>
</div>