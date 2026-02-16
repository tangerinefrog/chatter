import { wsStore } from "./store";
import type { WsEvent } from "$lib/websocket/event"

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
            
            break;
        default:
            console.warn('unknown event type:', (event as any).type);
    }
}