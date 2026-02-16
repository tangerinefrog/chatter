export type WsEvent = {
    type: 'send_message' | 'new_message';
    chat_id?: number;
    content?: string;
    sender_id?: number;
    message_id?: number;
    from_me?: boolean;
    date?: Date;
};
