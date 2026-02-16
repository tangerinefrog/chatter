import { wsStore } from "$lib/stores/websocket";
import { appendMessage } from '$lib/stores/messages';
import type { WsEvent } from "$lib/websocket/event"
import type { Message } from "$lib/models/message";

const socketUrl = 'ws://localhost:8080/ws';

let socket: WebSocket | null = null;
let isIntentionalClose = false;

export function connect() {
    wsStore.update(s => ({ ...s, status: 'connecting' }));
    
    socket = new WebSocket(socketUrl);
    socket.onopen = () => {
        wsStore.update(s => ({ ...s, status: 'connected' }));
    };

    socket.onclose = () => {
        wsStore.update(s => ({ ...s, socket: null, status: 'disconnected' }));
        if (!isIntentionalClose) {
            setTimeout(connect, 5000);
        }
    };

    socket.onerror = () => {
        socket?.close();
    };

    socket.onmessage = (event: MessageEvent) => {
        const message = JSON.parse(event.data) as WsEvent;
        handleEvent(message);
    };
}

export function disconnect() {
    isIntentionalClose = true;
    socket?.close;
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
                    createdAt: event.date?.toISOString() ?? ""
                }
                appendMessage(event.chat_id, message);
            }
            
            break;
        default:
            console.warn('unknown event type:', (event as any).type);
    }
}