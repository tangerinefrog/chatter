import type { Chat } from '$lib/models/chat';

const API_URL = 'http://localhost:8080/api';

export async function getChats(): Promise<Chat[]> {
    const res = await fetch(`${API_URL}/chats`, 
        {
            method: 'GET',
            credentials: 'include'
        }
    );

    const data = await res.json(); 
    if (!res.ok) {
        throw new Error(data?.message || 'Failed to get chats from server');
    }

    const chats: Chat[] = data.chats.map((chat :any) => {
        return {
            id: chat.id,
            type: chat.type as 'direct' | 'group',
            name: chat.name ?? null,
            lastMessage: chat.last_message ?? null,
            createdAt: chat.created_at
        };
    });
 
    return chats;
}