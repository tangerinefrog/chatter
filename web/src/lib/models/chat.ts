export type Chat = {
    id: string
    type: 'direct' | 'group'
    name: string | null
    lastMessage: string | null
    createdAt: string
};