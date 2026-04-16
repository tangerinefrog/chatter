export type Message = {
    id: number;
    fromMe: boolean;
    text: string;
    createdAt: Date;
    readAt: Date | null;
};