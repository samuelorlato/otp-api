package ports

import (
	"github.com/samuelorlato/otp-api.git/internal/handlers/dtos"
	"github.com/samuelorlato/otp-api.git/pkg/errors"
)

type OTPUsecase interface {
	SendOTP(emailDTO *dtos.Email) (*errors.HTTPError, *string)
	VerifyOTP(OTPVerificationDTO dtos.OTPVerification) (*errors.HTTPError)
}
