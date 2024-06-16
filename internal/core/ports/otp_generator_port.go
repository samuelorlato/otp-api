package ports

type OTPGenerator interface {
	Generate(length int, useLetters bool, useUpperCase bool, useSpecialChars bool) string
}
