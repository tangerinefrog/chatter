<script lang="ts">
    import { goto } from '$app/navigation';
    import AuthForm from '$lib/components/AuthForm.svelte';
    import { login, signup, type ApiError } from '$lib/api/client';

    let error = '';
    let isLoading = false;

    async function handleSubmit({ username, password, mode }: { username: string; password: string; mode: 'login' | 'signup' }) {
        error = '';
        isLoading = true;

        try {
            if (mode === 'login') {
                await login({ username, password });
            } else {
                await signup({ username, password });
            }
        } catch (err) {
            const apiError = err as ApiError;
            error = apiError.message || `${mode} failed`;
            isLoading = false;
            return; 
        }

        goto('/');
    }

    function handleToggleMode() {
        error = '';
    }
</script>

<AuthForm {error} {isLoading} onSubmit={handleSubmit} onToggleMode={handleToggleMode} />
