export type WsEvent = {
    type: 'send_message' | 'new_message' | 'read_message';
    chat_id?: string;
    content?: string;
    file_ids?: string[];
    files?: FileInfo[];
    sender_id?: string;
    message_id?: string;
    from_me?: boolean;
    date?: Date;
    read_at?: string | Date;
};

export type FileInfo = {
    id: string;
    name: string;
    mime_type: string;
    size_bytes: number;
};
