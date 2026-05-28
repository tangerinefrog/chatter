<script lang="ts">
    import { KEYBOARD_KEYS } from '$lib/constants';
    import paperclipIcon from '$lib/assets/icons/paperclip.svg?url';

    export let value = '';
    export let disabled = false;
    export let attachedFile: File | null = null;
    export let onSend: (file?: File | null) => void;

    let inputRef: HTMLInputElement;
    let fileInputRef: HTMLInputElement;

    export function focusInput() {
        inputRef?.focus();
    }

    function handleKeydown(event: KeyboardEvent) {
        if (event.key === KEYBOARD_KEYS.ENTER && !event.shiftKey) {
            event.preventDefault();
            onSend(attachedFile);
        }
    }

    function handleFileChange(event: Event) {
        const files = (event.target as HTMLInputElement)?.files;
        attachedFile = files?.[0] ?? null;
    }

    function clearAttachment() {
        attachedFile = null;
        if (fileInputRef) {
            fileInputRef.value = '';
        }
    }
</script>

<div class="chat-input">
    <label class="attachment-button" aria-label="Attach file">
        <input
            bind:this={fileInputRef}
            type="file"
            hidden
            on:change={handleFileChange}
            disabled={disabled}
        />
        <img src={paperclipIcon} alt="Attach file" />
    </label>

    <div class="input-wrapper">
        <input
            bind:this={inputRef}
            type="text"
            placeholder="Type a message..."
            bind:value={value}
            on:keydown={handleKeydown}
            aria-label="Message input"
            disabled={disabled}
        />

        {#if attachedFile}
            <div class="attachment-preview">
                <span>{attachedFile.name}</span>
                <button type="button" on:click={clearAttachment} aria-label="Remove attachment">
                    ×
                </button>
            </div>
        {/if}
    </div>

    <button on:click={() => onSend(attachedFile)} aria-label="Send message" disabled={disabled || (!value.trim() && !attachedFile)}>
        Send
    </button>
</div>
