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