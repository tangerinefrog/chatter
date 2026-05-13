<script lang="ts">
    import type { Chat } from '$lib/models/chat';
    import type { Message } from '$lib/models/message';
    import { getChats, getMessages, createChat } from '$lib/api/client';
    import { connect, disconnect, sendEvent, setOnNewMessageCallback } from '$lib/websocket/client';
    import { onMount, tick } from 'svelte';
    import { setMessages, messagesStore } from '$lib/stores/messages';
    import { KEYBOARD_KEYS } from '$lib/constants';
    import { formatDateShort, isSameDay } from '$lib/utils/date';
    import type { ApiError } from '$lib/api/client';
    import ChatList from '$lib/components/ChatList.svelte';
    import ChatInput from '$lib/components/ChatInput.svelte';
    import MessageRow from '$lib/components/MessageRow.svelte';
    import NewChatModal from '$lib/components/NewChatModal.svelte';
    
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
    let chatInput: { focusInput?: () => void } | null = null;
    let chatPages: Record<string, number> = {};
    let chatHasMore: Record<string, boolean> = {};

    $: messages = currentChat ? $messagesStore[currentChat.id] ?? [] : [];

    async function selectChat(chat: Chat) {
        tick().then(() => {
            chatInput?.focusInput?.();
        });

        if (currentChat?.id === chat.id) {
            return;
        }

        currentChat = chat;
        chatPages[chat.id] = 1;
        chatHasMore[chat.id] = true;

        await loadMessages(chat.id, 1);
        await markVisibleMessagesAsRead();
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
                createdAt: chat.created_at,
                unreadMessagesCount: chat.unread_messages_count,
            }));
        } catch (err) {
            const apiError = err as ApiError;
            error = apiError.message || 'Failed to load chats';
            chats = [];
        } finally {
            isLoadingChats = false;
        }
    }

    async function loadMessages(chatID: string, page: number, prepend?: boolean) {
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
                createdAt: new Date(message.created_at),
                readAt: message.read_at ? new Date(message.read_at) : null
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

    function isRowVisible(row: HTMLElement): boolean {
        if (!messagesContainer) {
            return false;
        }

        const containerRect = messagesContainer.getBoundingClientRect();
        const rect = row.getBoundingClientRect();

        return rect.bottom > containerRect.top && rect.top < containerRect.bottom;
    }

    async function markVisibleMessagesAsRead() {
        if (!currentChat || !messagesContainer || messages.length === 0) {
            return;
        }

        if (document.visibilityState !== 'visible' || !document.hasFocus()) {
            return;
        }

        const unreadRows = Array.from(
            messagesContainer.querySelectorAll('.message-row:not(.me)[data-read="false"]')
        ) as HTMLElement[];

        const visibleUnread = unreadRows.filter(isRowVisible);

        if (visibleUnread.length === 0) {
            return;
        }

        const lastVisible = visibleUnread.reduce((prev, row) => {
            const prevDate = new Date(prev.dataset.messageDate ?? '');
            const rowDate = new Date(row.dataset.messageDate ?? '');
            return rowDate > prevDate ? row : prev;
        });

        const messageID = lastVisible.dataset.messageId;
        if (messageID) {
            await markAsRead(messageID);
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
            await markVisibleMessagesAsRead();
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

        await markVisibleMessagesAsRead();
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
        error = null;
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

    async function handleNewMessage(chatId: string, message: Message) {
        const chatIndex = chats.findIndex((chat) => chat.id === chatId);

        if (chatIndex === -1) {
            refreshChats();
            return;
        }

        const unreadMessagesCount = chats[chatIndex].unreadMessagesCount ?? 0;

        const updatedChat: Chat = {
            ...chats[chatIndex],
            lastMessage: message.text,
            lastMessageDate: message.createdAt,
            unreadMessagesCount: !message.fromMe ? unreadMessagesCount + 1 : 0
        };

        chats = [
            updatedChat,
            ...chats.filter(chat => chat.id !== chatId)
        ];

        if (currentChat?.id === chatId) {
            currentChat = { ...updatedChat };
            if (message.fromMe) {
                await scrollToBottom();
            } else if (document.visibilityState === 'visible' && document.hasFocus()) {
                await markVisibleMessagesAsRead();
            }
        }
    }

    async function markAsRead(messageID: string | undefined) {
        if (!currentChat?.id) {
            return;
        }

        const chatId = currentChat.id;
        chats = chats.map((chat) =>
            chat.id === chatId ? { ...chat, unreadMessagesCount: 0 } : chat
        );
        currentChat = currentChat.id === chatId ? { ...currentChat, unreadMessagesCount: 0 } : currentChat;
        sendEvent({
            type: 'read_message',
            chat_id: chatId,
            message_id: messageID
        });
    }

    function handleVisibilityChange() {
        if (document.visibilityState === 'visible') {
            markVisibleMessagesAsRead();
        }
    }

    onMount(() => {
        connect();
        refreshChats();
        setOnNewMessageCallback(handleNewMessage);

        document.addEventListener('visibilitychange', handleVisibilityChange);

        return () => {
            disconnect();
            document.removeEventListener('visibilitychange', handleVisibilityChange);
        };
    });
    
</script>

<svelte:window on:keydown={handleKeydown} />

<div class="app">
    <ChatList
        {chats}
        currentChatId={currentChat?.id ?? null}
        isLoading={isLoadingChats}
        onSelectChat={selectChat}
        onAddNewChat={openModal}
    />

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
                    {#each messages as msg, index}
                    {#if index === 0 || !isSameDay(messages[index - 1].createdAt, msg.createdAt)}
                        <div class="message-date-divider">{formatDateShort(msg.createdAt)}</div>
                    {/if}
                    <MessageRow message={msg} />
                {/each}
            {/if}
            </div>

            <ChatInput
                bind:this={chatInput}
                bind:value={currentMessage}
                disabled={!currentChat}
                onSend={sendMessage}
            />
        {:else}
            <div class="empty">Select a chat or start a new one</div>
        {/if}
    </main>

    {#if isNewChatModalOpen}
        <NewChatModal
            {error}
            bind:username={username}
            onSubmit={addChat}
            onClose={closeModal}
        />
    {/if}
</div>
