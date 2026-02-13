package dto

type CreateChatRequest struct {
	IsDirect             bool     `json:"is_direct" validate:"required"`
	Name                 string   `json:"name" validate:"max=64"`
	ParticipantUsernames []string `json:"participant_usernames" validate:"required"`
}

type NewChatResponse struct {
	ID int32 `json:"id"`
}
