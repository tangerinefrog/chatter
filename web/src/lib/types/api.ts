export interface ApiChat {
    id: string;
    type: 'direct' | 'group';
    name: string | null;
    last_message: string | null;
    last_message_date: string | null;
    created_at: string;
    unread_messages_count: number | null;
}

export interface ApiMessage {
    id: string;
    content: string;
    from_me: boolean;
    user_id: string;
    created_at: string;
    read_at: string | null;
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
