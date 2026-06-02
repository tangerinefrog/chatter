export type Message = {
    id: string;
    fromMe: boolean;
    text: string;
    createdAt: Date;
    readAt: Date | null;
    files?: MessageFile[];
};

export type MessageFile = {
    id: string;
    name: string;
    mimeType: string;
    sizeBytes: number;
};