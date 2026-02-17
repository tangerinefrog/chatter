export interface ApiChat {
    id: number;
    type: 'direct' | 'group';
    name: string | null;
    last_message: string | null;
    last_message_date: string | null;
    created_at: string;
}

export interface ApiMessage {
    id: number;
    content: string;
    from_me: boolean;
    user_id: number;
    created_at: string;
}

export interface ChatsResponse {
    chats: ApiChat[];
}

export interface MessagesResponse {
    messages: ApiMessage[];
}

export interface CreateChatRequest {
    is_direct: boolean;
    participant_usernames: string[];
}

export interface AuthRequest {
    username: string;
    password: string;
}

export interface AuthResponse {
    message?: string;
}
