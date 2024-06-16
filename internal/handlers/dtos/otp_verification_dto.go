package dtos

type OTPVerification struct {
	OTP                   string `json:"OTP" binding:"required"`
	EncryptedVerification string `json:"encryptedVerification" binding:"required"`
	Email                 string `json:"email" binding:"required"`
}
