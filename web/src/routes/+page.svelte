<script lang="ts">
    import type { Chat } from '$lib/models/chat';
    import type { Message } from '$lib/models/message';
    import { apiFetch } from '$lib/api/client';
    import { connect, disconnect, sendEvent } from '$lib/websocket/client';
    import { onMount } from 'svelte';
    
    import '$lib/css/main.css';
    
    let chats: Chat[] = [];
    let messages: Message[] = [];
    let currentChat: Chat | null = null;
    let currentMessage = '';
    let isNewChatModalOpen = false;
    let username = '';

    function selectChat(chat: Chat) {
        if (currentChat?.id === chat.id) {
            return;
        }

        currentChat = chat;   
        loadMessages(chat.id, 1);
    }

    async function refreshChats() {
        try {
            const resp = await apiFetch('/chats', {
                method: 'GET'
            });

            chats = resp.chats.map((chat :any) => {
                return {
                    id: chat.id,
                    type: chat.type as 'direct' | 'group',
                    name: chat.name ?? null,
                    lastMessage: chat.last_message ?? null,
                    createdAt: chat.created_at
                };
            });
        } catch (err: any) {
            console.warn('Could not load chats from backend:', err);
            chats = [];
        }
    }

    async function loadMessages(chatID: number, page: number) {
        try {
            const resp = await apiFetch(`/chats/${chatID}/messages?page=${page}`, {
                method: 'GET'
            });

            messages = resp.messages.map((message :any) => {
                return {
                    id: message.id,
                    text: message.content,
                    fromMe: message.from_me,
                    userId: message.user_id,
                    createdAt: message.created_at
                };
            });
        } catch (err: any) {
            console.warn('Could not load messages from backend:', err);
        }
    }

    async function sendMessage() {
        const messageClean = currentMessage.trim();

        if (!messageClean || !currentChat?.id) {
            return;
        }        

        sendEvent({
            type: 'send_message',
            chat_id: currentChat.id,
            content: messageClean
        });

        currentMessage = '';
    }    

    function openModal() {
        isNewChatModalOpen = true;
    }

    function closeModal() {
        isNewChatModalOpen = false;
        username = '';
    }

    async function addChat() {
        if (!username.trim()) {
            return;
        }
        const req = {
            is_direct: true,
            participant_usernames: [username]
        }

        try {
            await apiFetch('/chats', {
                method: 'POST',
                credentials: 'include',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify(req),
            });

            await refreshChats();
        } catch (err: any) {
            console.warn('Could not create chat:', err);
        }

        closeModal();
    }

    function handleKeydown(event: KeyboardEvent) {
        if (event.key === 'Escape') {
            if (isNewChatModalOpen) {
                closeModal();
            } else {
                currentChat = null;
            }
        }
    }

    onMount(() => {
        connect();
        refreshChats();
    });
    
</script>

<svelte:window on:keydown={handleKeydown} />

<div class="app">
    <aside class="sidebar">
        <div class="sidebar-header">
            Chats

            <button class="button" on:click={openModal}>
                + Add new
            </button>
        </div>

        <div class="contacts">
            {#each chats as chat}
                <div role="button" tabindex="0" class="contact {chat.id === currentChat?.id ? 'active' : ''}"
                    on:click={() => selectChat(chat)}
                    on:keydown={() => selectChat(chat)}                    
                >
                    <div class="contact-name">{chat.name}</div>
                    <div class="contact-last">{chat.lastMessage}</div>
                </div>
            {/each}
        </div>
    </aside>

    <main class="chat">
        {#if currentChat}
            <div class="chat-header">
                {currentChat.name}
            </div>

            <div class="messages">
                {#each messages as msg}
                    <div class="message-row {msg.fromMe ? 'me' : 'them'}">
                        <div class="bubble">
                            {msg.text}
                        </div>
                    </div>
                {/each}
            </div>

            <div class="chat-input">
                <input
                    type="text"
                    placeholder="Type a message..."
                    bind:value={currentMessage}
                    on:keydown={(e) => e.key === 'Enter' && sendMessage()}
                />
                <button on:click={sendMessage}>
                    Send
                </button>
            </div>
        {:else}
            <div class="empty">
                Select a chat or start a new one
            </div>
        {/if}
    </main>

    {#if isNewChatModalOpen}
    <div class="backdrop">
        <div class="modal">
        <h2>Add new chat</h2>

        <input
            class="input"
            type="text"
            placeholder="Username"
            bind:value={username}
        />

        <div class="actions">
            <button class="button" on:click={addChat}>Add</button>
            <button class="button button-secondary" on:click={closeModal}>Cancel</button>
        </div>
        </div>
    </div>
    {/if}
</div>
