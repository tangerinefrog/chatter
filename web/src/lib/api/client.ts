import type { ChatsResponse, MessagesResponse, CreateChatRequest, AuthRequest, AuthResponse } from '$lib/types/api';

const API_URL = 'http://localhost:8080/api';

export interface ApiError extends Error {
    status?: number;
}

function createApiError(message: string, status?: number): ApiError {
    const error = new Error(message) as ApiError;
    error.status = status;
    return error;
}

function extractErrorMessage(error: unknown): string {
    if (error instanceof Error) return error.message;
    if (error && typeof error === 'object' && 'message' in error) {
        return String(error.message);
    }
    return String(error || 'Network error occurred');
}

export async function apiFetch<T = unknown>(
    path: string, 
    options: RequestInit = {}
): Promise<T> {
    const url = `${API_URL}${path}`;
    const headers = {
        'Content-Type': 'application/json',
        ...options.headers
    };

    try {
        const res = await fetch(url, {
            credentials: 'include',
            headers,
            ...options
        });

        const data = await res.json().catch(() => null);

        if (!res.ok) {
            throw createApiError(
                data?.message || `HTTP ${res.status}`,
                res.status
            );
        }

        return data as T;
    } catch (error) {
        if (error instanceof Error && 'status' in error) {
            throw error;
        }
        throw new Error(extractErrorMessage(error));
    }
}

function postJson<T>(path: string, body: unknown): Promise<T> {
    return apiFetch<T>(path, {
        method: 'POST',
        body: JSON.stringify(body)
    });
}

export async function getChats(): Promise<ChatsResponse> {
    return apiFetch<ChatsResponse>('/chats');
}

export async function getMessages(chatId: number, page: number): Promise<MessagesResponse> {
    return apiFetch<MessagesResponse>(`/chats/${chatId}/messages?page=${page}`);
}

export async function createChat(request: CreateChatRequest): Promise<void> {
    await postJson('/chats', request);
}

export async function login(request: AuthRequest): Promise<AuthResponse> {
    return postJson<AuthResponse>('/auth/login', request);
}

export async function signup(request: AuthRequest): Promise<AuthResponse> {
    return postJson<AuthResponse>('/auth/signup', request);
}