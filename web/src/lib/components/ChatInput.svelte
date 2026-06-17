<script lang="ts">
    import { KEYBOARD_KEYS } from "$lib/constants";
    import paperclipIcon from "$lib/assets/icons/paperclip.svg?url";
    import emojiIcon from "$lib/assets/icons/emoji.svg?url";
    import { onMount } from "svelte";

    export let value = "";
    export let disabled = false;
    export let attachedFile: File | null = null;
    export let onSend: (file?: File | null) => void;

    let inputRef: HTMLInputElement;
    let fileInputRef: HTMLInputElement;
    let showEmojiPicker = false;

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
            fileInputRef.value = "";
        }
    }

    function handleEmojiClick(event: any) {
        value += event.detail.unicode;
        inputRef?.focus();
    }

    function handlePickerClickOutside(event: any) {
        if (!event.target.closest(".emoji-picker-wrapper")) {
            showEmojiPicker = false;
        }
    }

    onMount(async () => {
        await import("emoji-picker-element");
    });
</script>

<svelte:window on:click={handlePickerClickOutside} />

<div class="chat-input">
    <label class="attachment-button" aria-label="Attach file">
        <input
            bind:this={fileInputRef}
            type="file"
            hidden
            on:change={handleFileChange}
            {disabled}
        />
        <img src={paperclipIcon} alt="Attach file" />
    </label>

    <div class="input-wrapper">
        <div class="input-row">
            <div class="input-with-emoji">
                <input
                    bind:this={inputRef}
                    type="text"
                    placeholder="Type a message..."
                    bind:value
                    on:keydown={handleKeydown}
                    aria-label="Message input"
                    {disabled}
                />
                <div class="emoji-picker-wrapper">
                    <button
                        type="button"
                        class="emoji-button"
                        on:click|stopPropagation={() =>
                            (showEmojiPicker = !showEmojiPicker)}
                        aria-label="Open emoji picker"
                        {disabled}
                    >
                        <img src={emojiIcon} alt="Emoji" />
                    </button>
                    {#if showEmojiPicker}
                        <div
                            class="emoji-picker-popover"
                            on:click|stopPropagation
                        >
                            <emoji-picker
                                class="dark"
                                on:emoji-click={handleEmojiClick}
                            />
                        </div>
                    {/if}
                </div>
            </div>
        </div>

        {#if attachedFile}
            <div class="attachment-preview">
                <span>{attachedFile.name}</span>
                <button
                    type="button"
                    on:click={clearAttachment}
                    aria-label="Remove attachment">×</button
                >
            </div>
        {/if}
    </div>

    <button
        on:click={() => onSend(attachedFile)}
        aria-label="Send message"
        disabled={disabled || (!value.trim() && !attachedFile)}
    >
        Send
    </button>
</div>
