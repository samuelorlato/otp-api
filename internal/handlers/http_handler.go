package handlers

import (
	"net/mail"

	"github.com/gin-gonic/gin"
	"github.com/samuelorlato/otp-api.git/internal/core/ports"
	"github.com/samuelorlato/otp-api.git/internal/handlers/dtos"
	"github.com/samuelorlato/otp-api.git/pkg/errors"
)

type HTTPHandler struct {
	engine       *gin.Engine
	OTPUsecase   ports.OTPUsecase
	errorHandler *ErrorHandler
}

func NewHTTPHandler(engine *gin.Engine, OTPUsecase ports.OTPUsecase, errorHandler *ErrorHandler) *HTTPHandler {
	return &HTTPHandler{
		engine:       engine,
		OTPUsecase:   OTPUsecase,
		errorHandler: errorHandler,
	}
}

func (h *HTTPHandler) SetRoutes() {
	h.engine.POST("/email/otp", h.sendOTPToEmail)
	h.engine.POST("/verify/otp", h.verifyOTP)
}

func (h *HTTPHandler) sendOTPToEmail(c *gin.Context) {
	var emailDTO dtos.Email

	bindErr := c.BindJSON(&emailDTO)
	if bindErr != nil {
		err := errors.NewValidationError(bindErr)
		h.errorHandler.Handle(err, c)
		return
	}

	_, parseErr := mail.ParseAddress(emailDTO.Email)
	if parseErr != nil {
		err := errors.NewValidationError(parseErr)
		h.errorHandler.Handle(err, c)
		return
	}

	err, details := h.OTPUsecase.SendOTP(&emailDTO)
	if err != nil {
		h.errorHandler.Handle(err, c)
		return
	}

	c.JSON(200, gin.H{"details": details})
}

func (h *HTTPHandler) verifyOTP(c *gin.Context) {
	var OTPVerificationDTO dtos.OTPVerification

	bindErr := c.BindJSON(&OTPVerificationDTO)
	if bindErr != nil {
		err := errors.NewValidationError(bindErr)
		h.errorHandler.Handle(err, c)
		return
	}

	_, parseErr := mail.ParseAddress(OTPVerificationDTO.Email)
	if parseErr != nil {
		err := errors.NewValidationError(parseErr)
		h.errorHandler.Handle(err, c)
		return
	}

	err := h.OTPUsecase.VerifyOTP(OTPVerificationDTO)
	if err != nil {
		h.errorHandler.Handle(err, c)
		return
	}

	c.JSON(200, gin.H{"details": "OTP verified", "email": OTPVerificationDTO.Email})
}
