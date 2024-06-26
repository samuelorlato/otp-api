package ports

import "github.com/samuelorlato/otp-api.git/pkg/errors"

type EmailSender interface {
	Send(to string, message string) *errors.HTTPError
}
