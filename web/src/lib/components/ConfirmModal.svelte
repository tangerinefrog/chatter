<script lang="ts">
    import { KEYBOARD_KEYS } from '$lib/constants';

    export let message: string;
    export let onConfirm: () => void;
    export let onClose: () => void;

    function handleKeydown(event: KeyboardEvent) {
        if (event.key === KEYBOARD_KEYS.ESCAPE) {
            event.preventDefault();
            onClose();
        } else if (event.key === KEYBOARD_KEYS.ENTER) {
            event.preventDefault();
            onConfirm();
        }
    }
</script>

<svelte:window on:keydown={handleKeydown} />

<div class="backdrop" role="dialog" aria-modal="true">
    <div class="modal">
        <p class="confirm-message">{message}</p>

        <div class="actions">
            <button class="button" on:click={onConfirm}>Yes</button>
            <button class="button button-secondary" on:click={onClose}>Cancel</button>
        </div>
    </div>
</div>

<style>
    .confirm-message {
        margin: 0 0 1.5rem;
        color: var(--text-main);
    }
</style>
