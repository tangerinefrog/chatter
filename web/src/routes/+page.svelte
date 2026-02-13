<script lang="ts">
    import '$lib/css/main.css';
    import type { Chat } from '$lib/models/chat';
    import type { Message } from '$lib/models/message';

    let chats: Chat[] = [
        { 
            id: "dfgkj1kfdsj",
            type: 'direct',
            name: 'someone',
            createdAt: new Date().toISOString(),
            lastMessage: 'yo my guy'
        },{ 
            id: "fdsqweqwe",
            type: 'direct',
            name: 'someone',
            createdAt: new Date().toISOString(),
            lastMessage: 'yo my guy'
        },
    ];

    let currentChat = chats[0];
    let currentMessage = '';

    let messages: Message[] = [
        { id: 'qfrjh52we', fromMe: false, text: 'Hey there!' },
        { id: 'qgsjhq2we', fromMe: true, text: 'Hi! What\'s up?' },
        { id: 'qfs1hq2ge', fromMe: false, text: 'Just testing this chat UI.' },
        { id: 'qfsjghq2we', fromMe: true, text: 'Looks decent so far.' }
    ];

    function selectChat(chat: Chat) {
        if (currentChat.id === chat.id) {
            return;
        }

        currentChat = chat;   
        messages = [];

    }

    function sendMessage() {
        if (!currentMessage.trim()) return;

        messages = [
            ...messages,
            { id: "qweasd", fromMe: true, text: currentMessage }
        ];

        currentMessage = '';
    }
</script>

<div class="app">
    <aside class="sidebar">
        <div class="sidebar-header">
            Chats
        </div>

        <div class="contacts">
            {#each chats as chat}
                <div role="button" tabindex="0" class="contact {chat.id === currentChat.id ? 'active' : ''}"
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
                Select a conversation
            </div>
        {/if}
    </main>
</div>
