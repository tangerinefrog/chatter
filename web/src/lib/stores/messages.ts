import type { Message } from '$lib/models/message';
import { writable } from 'svelte/store';

interface MessagesStore {
    [chatId: number]: Message[];
}

export const messagesStore = writable<MessagesStore>({});

export function appendMessage(chatId: number, message: Message): void {
    messagesStore.update(store => ({
        ...store,
        [chatId]: [...(store[chatId] ?? []), message]
    }));
}

export function setMessages(chatId: number, messages: Message[]): void {
    messagesStore.update(store => ({
        ...store,
        [chatId]: messages
    }));
}

export function markSentMessagesAsRead(chatId: number, lastMessageId: number, readAt: Date): void {
    messagesStore.update(store => {
        const msgs = store[chatId];
        if (!msgs) return store;

        return {
            ...store,
            [chatId]: msgs.map(msg =>
                msg.id <= lastMessageId && msg.fromMe && !msg.readAt
                    ? { ...msg, readAt }
                    : msg
            )
        };
    });
}