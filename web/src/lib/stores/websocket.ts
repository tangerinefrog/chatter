import { writable } from 'svelte/store';

interface WsStore {
    socket: WebSocket | null;
    status: 'disconnected' | 'connecting' | 'connected';
}

export const wsStore = writable<WsStore>({
    socket: null,
    status: 'disconnected'
});