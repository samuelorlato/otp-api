package models

import "time"

type OTPDetails struct {
	ID        string
	Timestamp time.Time
	Email     string
	Success   bool
	Message   string
}
