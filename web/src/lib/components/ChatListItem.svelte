<script lang="ts">
    import type { Chat } from '$lib/models/chat';
    import { formatTimestamp, formatDateShort, isToday } from '$lib/utils/date';

    export let chat: Chat;
    export let selected = false;
    export let onSelect: (chat: Chat) => void;

    function handleClick() {
        onSelect(chat);
    }

    function handleKeydown(event: KeyboardEvent) {
        if (event.key === 'Enter' || event.key === ' ') {
            event.preventDefault();
            onSelect(chat);
        }
    }
</script>

<div
    role="button"
    tabindex="0"
    class="contact {selected ? 'active' : ''}"
    on:click={handleClick}
    on:keydown={handleKeydown}
    aria-label="Chat with {chat.name || 'Unknown'}"
>
    <div class="contact-header">
        <div class="contact-name">{chat.name || 'Unknown'}</div>
        {#if chat.unreadMessagesCount && chat.unreadMessagesCount > 0}
            <div class="unread-badge">{chat.unreadMessagesCount}</div>
        {:else if chat.lastMessageDate}
            <div class="contact-date">
                {#if isToday(chat.lastMessageDate)}
                    {formatTimestamp(chat.lastMessageDate)}
                {:else}
                    {formatDateShort(chat.lastMessageDate)}
                {/if}
            </div>
        {/if}
    </div>
    <div class="contact-last">{chat.lastMessage || 'No messages'}</div>
</div>
