<script lang="ts">
    import { KEYBOARD_KEYS } from '$lib/constants';

    export let value = '';
    export let disabled = false;
    export let onSend: () => void;

    let inputRef: HTMLInputElement;

    export function focusInput() {
        inputRef?.focus();
    }

    function handleKeydown(event: KeyboardEvent) {
        if (event.key === KEYBOARD_KEYS.ENTER && !event.shiftKey) {
            event.preventDefault();
            onSend();
        }
    }
</script>

<div class="chat-input">
    <input
        bind:this={inputRef}
        type="text"
        placeholder="Type a message..."
        bind:value={value}
        on:keydown={handleKeydown}
        aria-label="Message input"
        disabled={disabled}
    />
    <button on:click={onSend} aria-label="Send message" disabled={disabled || !value.trim()}>
        Send
    </button>
</div>
