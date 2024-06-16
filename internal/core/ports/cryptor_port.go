package ports

import "github.com/samuelorlato/otp-api.git/pkg/errors"

type Cryptor interface {
	Encrypt(text string) (*errors.HTTPError, *string)
	Decrypt(text string) (*errors.HTTPError, *string)
}
