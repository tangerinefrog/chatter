<script lang="ts">
    import { KEYBOARD_KEYS } from '$lib/constants';

    export let error: string | null = null;
    export let username = '';
    export let onSubmit: () => void;
    export let onClose: () => void;

    function handleKeydown(event: KeyboardEvent) {
        if (event.key === KEYBOARD_KEYS.ESCAPE) {
            event.preventDefault();
            onClose();
        } else if (event.key === KEYBOARD_KEYS.ENTER) {
            event.preventDefault();
            onSubmit();
        }
    }
</script>

<div class="backdrop" role="dialog" aria-modal="true" aria-labelledby="modal-title">
    <div class="modal">
        <h2 id="modal-title">Add new chat</h2>

        {#if error}
            <div class="error" role="alert">{error}</div>
        {/if}

        <input
            class="input"
            type="text"
            placeholder="Username"
            bind:value={username}
            on:keydown={handleKeydown}
            aria-label="Username"
        />

        <div class="actions">
            <button class="button" on:click={onSubmit} disabled={!username.trim()}>
                Add
            </button>
            <button class="button button-secondary" on:click={onClose}>
                Cancel
            </button>
        </div>
    </div>
</div>
