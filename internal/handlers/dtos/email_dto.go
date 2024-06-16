package dtos

type Email struct {
	Email string `json:"email" binding:"required"`
}
