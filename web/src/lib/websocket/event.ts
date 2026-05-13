export type WsEvent = {
    type: 'send_message' | 'new_message' | 'read_message';
    chat_id?: string;
    content?: string;
    sender_id?: string;
    message_id?: string;
    from_me?: boolean;
    date?: Date;
    read_at?: string | Date;
};
