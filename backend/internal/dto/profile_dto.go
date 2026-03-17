package dto

type CreateProfileRequest struct {
	Name      string `json:"name" binding:"required"`
	AvatarURL string `json:"avatar_url"`
}

type ProfileResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	CreatedAt string `json:"created_at"`
}
