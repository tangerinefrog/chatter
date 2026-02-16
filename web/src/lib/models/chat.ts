export type Chat = {
    id: number
    type: 'direct' | 'group'
    name: string | null
    lastMessage: string | null
    lastMessageDate: Date | null
    createdAt: string
};