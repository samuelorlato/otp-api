package otp

import "math/rand"

type OTPGenerator struct {
	numbers          string
	lowerCaseLetters string
	upperCaseLetters string
	specialChars     string
}

func NewOTPGenerator() *OTPGenerator {
	return &OTPGenerator{
		numbers:          "0123456789",
		lowerCaseLetters: "abcdefghijklmnopqrstuvwxyz",
		upperCaseLetters: "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		specialChars:     "@#&!",
	}
}

func (o *OTPGenerator) Generate(length int, useLetters bool, useUpperCase bool, useSpecialChars bool) string {
	var charSource string

	charSource += o.numbers

	if useLetters {
		charSource += o.lowerCaseLetters
	}
	if useUpperCase {
		charSource += o.upperCaseLetters
	}
	if useSpecialChars {
		charSource += o.specialChars
	}

	b := make([]byte, length)
	for i := 0; i < length; i++ {
		randomIndex := rand.Intn(len(charSource))
		b[i] = charSource[randomIndex]
	}

	return string(b[:])
}
