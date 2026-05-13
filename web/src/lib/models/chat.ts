export type Chat = {
    id: string
    type: 'direct' | 'group'
    name: string | null
    lastMessage: string | null
    lastMessageDate: Date | null
    createdAt: string
    unreadMessagesCount: number | null
};