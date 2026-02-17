<script lang="ts">
    import type { Chat } from '$lib/models/chat';
    import type { Message } from '$lib/models/message';
    import { getChats, getMessages, createChat } from '$lib/api/client';
    import { connect, disconnect, sendEvent, setOnNewMessageCallback } from '$lib/websocket/client';
    import { onMount, onDestroy, tick } from 'svelte';
    import { setMessages, messagesStore } from '$lib/stores/messages';
    import { formatTimestamp } from '$lib/utils/date';
    import type { ApiError } from '$lib/api/client';
    import { KEYBOARD_KEYS } from '$lib/constants';
    
    import '$lib/css/main.css';
    
    const MESSAGES_PAGE_SIZE = 20;
    const MESSAGES_LOAD_THRESHOLD = 15;    
    
    let chats: Chat[] = [];
    let currentChat: Chat | null = null;
    let currentMessage = '';
    let isNewChatModalOpen = false;
    let username = '';
    let isLoadingChats = false;
    let isLoadingMessages = false;
    let isLoadingMore = false;
    let error: string | null = null;
    let messagesContainer: HTMLDivElement;
    let messageInput: HTMLInputElement;
    let chatPages: Record<number, number> = {};
    let chatHasMore: Record<number, boolean> = {};

    $: messages = currentChat ? $messagesStore[currentChat?.id] ?? [] : [];

    function selectChat(chat: Chat) {
        tick().then(() => {
            messageInput.focus();
        });

        if (currentChat?.id === chat.id) {
            return;
        }

        currentChat = chat;
        chatPages[chat.id] = 1;
        chatHasMore[chat.id] = true;
        loadMessages(chat.id, 1);
    }

    async function refreshChats() {
        isLoadingChats = true;
        error = null;
        
        try {
            const resp = await getChats();
            chats = resp.chats.map((chat) => ({
                id: chat.id,
                type: chat.type as 'direct' | 'group',
                name: chat.name ?? null,
                lastMessage: chat.last_message ?? null,
                lastMessageDate: chat.last_message_date ? new Date(chat.last_message_date) : null,
                createdAt: chat.created_at
            }));
        } catch (err) {
            const apiError = err as ApiError;
            error = apiError.message || 'Failed to load chats';
            chats = [];
        } finally {
            isLoadingChats = false;
        }
    }

    async function loadMessages(chatID: number, page: number, prepend?: boolean) {
        if (!prepend) {
            isLoadingMessages = true;
        }

        error = null;
        
        const shouldPreserveScroll = prepend && currentChat?.id === chatID && !!messagesContainer;
        const previousScrollHeight = shouldPreserveScroll ? messagesContainer.scrollHeight : 0;

        try {
            const resp = await getMessages(chatID, page);
            const msgs = resp.messages.map((message) => ({
                id: message.id,
                text: message.content,
                fromMe: message.from_me,
                userId: message.user_id,
                createdAt: new Date(message.created_at)
            }));

            chatPages[chatID] = page;
            chatHasMore[chatID] = msgs.length === MESSAGES_PAGE_SIZE;

            if (prepend && currentChat?.id === chatID) {
                const existing = messages;
                setMessages(chatID, [...msgs, ...existing]);

                await tick();
                if (messagesContainer && shouldPreserveScroll) {
                    const newScrollHeight = messagesContainer.scrollHeight;
                    messagesContainer.scrollTop = newScrollHeight - previousScrollHeight;
                }
            } else {
                setMessages(chatID, msgs);

                if (msgs.length > 0) {
                    const lastMsg = msgs[msgs.length - 1];
                    const chatIndex = chats.findIndex(c => c.id === chatID);
                    if (chatIndex !== -1) {
                        chats = chats.map((chat, idx) => 
                            idx === chatIndex 
                                ? { ...chat, lastMessageDate: lastMsg.createdAt }
                                : chat
                        );
                        if (currentChat?.id === chatID) {
                            currentChat = { ...currentChat, lastMessageDate: lastMsg.createdAt };
                        }
                    }
                }

                await scrollToBottom();
            }
        } catch (err) {
            const apiError = err as ApiError;
            error = apiError.message || 'Failed to load messages';
        } finally {
            if (!prepend) {
                isLoadingMessages = false;
            }
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

    async function handleMessagesScroll() {
        if (!currentChat || !messagesContainer || isLoadingMore || isLoadingMessages) {
            return;
        }

        const currentChatId = currentChat.id;
        if (!chatHasMore[currentChatId]) {
            return;
        }

        if (messagesContainer.scrollTop <= MESSAGES_LOAD_THRESHOLD) {
            const nextPage = (chatPages[currentChatId] ?? 1) + 1;
            isLoadingMore = true;
            try {
                await loadMessages(currentChatId, nextPage, true);
            } finally {
                isLoadingMore = false;
            }
        }
    }

    async function scrollToBottom() {
        if (messagesContainer) {
            await tick();
            requestAnimationFrame(() => {
                messagesContainer.scrollTop = messagesContainer.scrollHeight;
            });
        }
    }    

    function openModal() {
        isNewChatModalOpen = true;
    }

    function closeModal() {
        isNewChatModalOpen = false;
        username = '';
    }

    async function addChat() {
        const trimmedUsername = username.trim();
        if (!trimmedUsername) {
            error = 'Username is required';
            return;
        }

        error = null;
        
        try {
            await createChat({
                is_direct: true,
                participant_usernames: [trimmedUsername]
            });

            await refreshChats();
            closeModal();
        } catch (err) {
            const apiError = err as ApiError;
            error = apiError.message || 'Failed to create chat';
        }
    }

    function handleKeydown(event: KeyboardEvent) {
        if (event.key === KEYBOARD_KEYS.ESCAPE) {
            if (isNewChatModalOpen) {
                closeModal();
            } else {
                currentChat = null;
            }
        }
    }

    function handleContactKeydown(event: KeyboardEvent, chat: Chat) {
        if (event.key === KEYBOARD_KEYS.ENTER || event.key === KEYBOARD_KEYS.SPACE) {
            event.preventDefault();
            selectChat(chat);
        }
    }

    async function handleNewMessage(chatId: number, message: Message) {
        const chatIndex = chats.findIndex(chat => chat.id === chatId);
        
        if (chatIndex === -1) {
            refreshChats();
            return;
        }

        const updatedChat: Chat = {
            ...chats[chatIndex],
            lastMessage: message.text,
            lastMessageDate: message.createdAt
        };

        chats = [
            updatedChat,
            ...chats.filter(chat => chat.id !== chatId)
        ];

        if (currentChat?.id === chatId) {
            currentChat = { ...updatedChat };
            if (message.fromMe) {
                await scrollToBottom();
            }
        }
    }

    onMount(() => {
        connect();
        refreshChats();
        setOnNewMessageCallback(handleNewMessage);
    });

    onDestroy(() => {
        disconnect();
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
            {#if isLoadingChats}
                <div class="loading">Loading chats...</div>
            {:else if chats.length === 0}
                <div class="empty-state">No chats yet. Create one to get started!</div>
            {:else}
                {#each chats as chat}
                    <div 
                        role="button" 
                        tabindex="0" 
                        class="contact {chat.id === currentChat?.id ? 'active' : ''}"
                        on:click={() => selectChat(chat)}
                        on:keydown={(e) => handleContactKeydown(e, chat)}
                        aria-label="Chat with {chat.name || 'Unknown'}"
                    >
                        <div class="contact-header">
                            <div class="contact-name">{chat.name || 'Unknown'}</div>
                            {#if chat.lastMessageDate}
                                <div class="contact-date">{formatTimestamp(chat.lastMessageDate)}</div>
                            {/if}
                        </div>
                        <div class="contact-last">{chat.lastMessage || 'No messages'}</div>
                    </div>
                {/each}
            {/if}
        </div>
    </aside>

    <main class="chat">
        {#if currentChat}
            <div class="chat-header">
                {currentChat.name}
            </div>

            <div class="messages" bind:this={messagesContainer} on:scroll={handleMessagesScroll}>
                {#if isLoadingMessages}
                    <div class="loading">Loading messages...</div>
                {:else if messages.length === 0}
                    <div class="empty-state">No messages yet. Start the conversation!</div>
                {:else}
                    {#each messages as msg}
                        <div class="message-row {msg.fromMe ? 'me' : 'them'}">
                            <div class="bubble">
                                <div class="message-text">{msg.text}</div>
                                <div class="message-timestamp" aria-label="Sent at {formatTimestamp(msg.createdAt)}">
                                    {formatTimestamp(msg.createdAt)}
                                </div>
                            </div>
                        </div>
                    {/each}
                {/if}
            </div>

            <div class="chat-input">
                <input
                    type="text"
                    placeholder="Type a message..."
                    bind:value={currentMessage}
                    bind:this={messageInput}
                    on:keydown={(e) => {
                        if (e.key === KEYBOARD_KEYS.ENTER && !e.shiftKey) {
                            e.preventDefault();
                            sendMessage();
                        }
                    }}
                    aria-label="Message input"
                    disabled={!currentChat}
                />
                <button 
                    on:click={sendMessage}
                    aria-label="Send message"
                    disabled={!currentChat || !currentMessage.trim()}
                >
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
                    on:keydown={(e) => {
                        if (e.key === 'Enter') {
                            e.preventDefault();
                            addChat();
                        } else if (e.key === KEYBOARD_KEYS.ESCAPE) {
                            e.preventDefault();
                            closeModal();
                        }
                    }}
                    aria-label="Username"
                />

                <div class="actions">
                    <button class="button" on:click={addChat} disabled={!username.trim()}>
                        Add
                    </button>
                    <button class="button button-secondary" on:click={closeModal}>
                        Cancel
                    </button>
                </div>
            </div>
        </div>
    {/if}
</div>
