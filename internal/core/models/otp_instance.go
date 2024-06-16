package models

import "time"

type OTPInstance struct {
	OTP            string
	ExpirationTime time.Time
	Verified       bool
}
