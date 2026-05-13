export type Message = {
    id: string;
    fromMe: boolean;
    text: string;
    createdAt: Date;
    readAt: Date | null;
};