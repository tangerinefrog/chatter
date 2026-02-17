import { appendMessage } from '$lib/stores/messages';
import type { WsEvent } from "$lib/websocket/event";
import type { Message } from "$lib/models/message";
import { WS_RECONNECT_DELAY } from '$lib/constants';

const socketUrl = import.meta.env.VITE_WS_URL || 'ws://localhost:8080/ws';

let socket: WebSocket | null = null;
let isIntentionalClose = false;
let onNewMessageCallback: ((chatId: number, message: Message) => void) | null = null;

export function setOnNewMessageCallback(callback: (chatId: number, message: Message) => void) {
    onNewMessageCallback = callback;
}

export function connect() {
    if (socket?.readyState === WebSocket.CONNECTING || socket?.readyState === WebSocket.OPEN) {
        return;
    }

    socket = new WebSocket(socketUrl);
    socket.onopen = () => { };

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
        console.error('ws not connected');
        return;
    }

    socket.send(JSON.stringify(event));
}

function handleEvent(event: WsEvent) {
    switch (event.type) {
        case 'new_message':
            if (event.chat_id && event.message_id) {
                const message: Message = {
                    fromMe: event.from_me ?? false,
                    id: event.message_id,
                    text: event.content ?? "",
                    createdAt: event.date instanceof Date 
                        ? event.date 
                        : event.date 
                            ? new Date(event.date)
                            : new Date()
                };
                appendMessage(event.chat_id, message);
                
                if (onNewMessageCallback) {
                    onNewMessageCallback(event.chat_id, message);
                }
            }
            break;
        default:
            console.warn('Unknown event type:', event.type);
    }
}