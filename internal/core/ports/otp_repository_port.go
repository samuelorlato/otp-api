package ports

import (
	"time"

	"github.com/samuelorlato/otp-api.git/internal/core/models"
	"github.com/samuelorlato/otp-api.git/pkg/errors"
)

type OTPRepository interface {
	Create(OTP string, expirationTime time.Time) (*errors.HTTPError, *string)
	Get(ID string) (*errors.HTTPError, *models.OTPInstance)
	Verify(ID string) *errors.HTTPError
	Delete(ID string) *errors.HTTPError
}
