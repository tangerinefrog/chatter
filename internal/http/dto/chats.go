package dto

type CreateChatRequest struct {
	IsDirect             bool     `json:"is_direct" validate:"required"`
	Name                 string   `json:"name" validate:"max=64"`
	ParticipantUsernames []string `json:"participant_usernames" validate:"required"`
}

type NewChatResponse struct {
	ID int32 `json:"id"`
}

type ListChatsForUserResponse struct {
	Chats []Chat `json:"chat"`
}

type Chat struct {
	ID           int32             `json:"id"`
	Type         string            `json:"type"`
	Name         string            `json:"name"`
	Participants []ChatParticipant `json:"participants"`
	LastMessage  string            `json:"last_message"`
}

type ChatParticipant struct {
	ID       int32  `json:"id"`
	Username string `json:"username"`
}
