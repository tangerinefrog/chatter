import { appendMessage, markSentMessagesAsRead } from '$lib/stores/messages';
import type { WsEvent } from "$lib/websocket/event";
import type { Message } from "$lib/models/message";
import { WS_RECONNECT_DELAY } from '$lib/constants';

const WS_URL = `ws://${import.meta.env.VITE_API_ADDR}/ws`;

let socket: WebSocket | null = null;
let isIntentionalClose = false;
let pendingEvents: WsEvent[] = [];
let onNewMessageCallback: ((chatId: string, message: Message) => void) | null = null;

export function setOnNewMessageCallback(callback: (chatId: string, message: Message) => void) {
    onNewMessageCallback = callback;
}

export function connect() {
    if (socket?.readyState === WebSocket.CONNECTING || socket?.readyState === WebSocket.OPEN) {
        return;
    }

    socket = new WebSocket(WS_URL);
    socket.onopen = () => {
        while (pendingEvents.length > 0 && socket?.readyState === WebSocket.OPEN) {
            const event = pendingEvents.shift();
            if (event) {
                socket.send(JSON.stringify(event));
            }
        }
    };

    socket.onclose = (event) => {
        if (!isIntentionalClose) {
            setTimeout(connect, WS_RECONNECT_DELAY);
        }
    };

    socket.onerror = (error) => {
        console.error('WebSocket error:', error);
    };

    socket.onmessage = (event: MessageEvent) => {
        try {
            const message = JSON.parse(event.data) as WsEvent;
            handleEvent(message);
        } catch (error) {
            console.error('Failed to parse WebSocket message:', error);
        }
    };
}

export function disconnect() {
    isIntentionalClose = true;
    if (socket) {
        socket.close();
        socket = null;
    }
}

export function sendEvent(event: WsEvent) {
    if (!socket || socket.readyState !== WebSocket.OPEN) {
        pendingEvents.push(event);
        if (!socket || socket.readyState === WebSocket.CLOSED) {
            connect();
        }
        return;
    }

    socket.send(JSON.stringify(event));
}

async function handleEvent(event: WsEvent) {
    switch (event.type) {
        case 'new_message':
            if (event.chat_id && event.message_id) {
                const files = event.files?.map(f => ({
                    id: f.id,
                    name: f.name,
                    mimeType: f.mime_type,
                    sizeBytes: f.size_bytes,
                })) ?? [];

                const message: Message = {
                    fromMe: event.from_me ?? false,
                    id: event.message_id,
                    text: event.content ?? "",
                    createdAt: event.date instanceof Date 
                        ? event.date 
                        : event.date 
                            ? new Date(event.date)
                            : new Date(),
                    readAt: null,
                    files: files.length > 0 ? files : undefined,
                };
                appendMessage(event.chat_id, message);
                
                if (onNewMessageCallback) {
                    onNewMessageCallback(event.chat_id, message);
                }
            }
            break;
        case 'read_message':
            if (event.chat_id && event.message_id) {
                const readAt = event.read_at ? new Date(event.read_at) : new Date();
                markSentMessagesAsRead(event.chat_id, event.message_id, readAt);
            }
            break;
        default:
            console.warn('Unknown event type:', event.type);
    }
}