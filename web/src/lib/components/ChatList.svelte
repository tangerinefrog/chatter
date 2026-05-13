<script lang="ts">
    import type { Chat } from '$lib/models/chat';
    import ChatListItem from './ChatListItem.svelte';

    export let chats: Chat[] = [];
    export let currentChatId: string | null = null;
    export let isLoading = false;
    export let onSelectChat: (chat: Chat) => void;
    export let onAddNewChat: () => void;
</script>

<aside class="sidebar">
    <div class="sidebar-header">
        Chats

        <button class="button" on:click={onAddNewChat}>+ Add new</button>
    </div>

    <div class="contacts">
        {#if isLoading}
            <div class="loading">Loading chats...</div>
        {:else if chats.length === 0}
            <div class="empty-state">No chats yet. Create one to get started!</div>
        {:else}
            {#each chats as chat}
                <ChatListItem
                    {chat}
                    selected={chat.id === currentChatId}
                    onSelect={onSelectChat}
                />
            {/each}
        {/if}
    </div>
</aside>
