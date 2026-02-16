import type { Chat } from '$lib/models/chat';

const API_URL = 'http://localhost:8080/api';

export async function apiFetch(path: string, options: RequestInit = {}): Promise<any> {
    const url = `${API_URL}${path}`

    const res = await fetch(url, {
        credentials: 'include',
        headers: {
        'Content-Type': 'application/json',
        ...options.headers
        },
        ...options
    });

    const data = await res.json().catch(() => null);

    if (!res.ok) {
        throw new Error(data?.message || `HTTP ${res.status}`);
    }

    return data;
}