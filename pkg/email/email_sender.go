package email

import (
	"net/smtp"

	"github.com/samuelorlato/otp-api.git/pkg/errors"
)

type EmailSender struct {
	host     string
	port     string
	from     string
	password string
}

func NewEmailSender(host string, port string, from string, password string) *EmailSender {
	return &EmailSender{
		host:     host,
		port:     port,
		from:     from,
		password: password,
	}
}

func (e *EmailSender) Send(to string, message string) *errors.HTTPError {
	auth := smtp.PlainAuth("", e.from, e.password, e.host)

	toList := []string{to}
	messageBytes := []byte(message)

	hostAndPort := e.host + ":" + e.port

	err := smtp.SendMail(hostAndPort, auth, e.from, toList, messageBytes)
	if err != nil {
		err := errors.NewGenericError(err)
		return err
	}

	return nil
}
