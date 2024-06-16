package usecases

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/fatih/structs"
	"github.com/samuelorlato/otp-api.git/internal/core/models"
	"github.com/samuelorlato/otp-api.git/internal/core/ports"
	"github.com/samuelorlato/otp-api.git/internal/handlers/dtos"
	"github.com/samuelorlato/otp-api.git/pkg/errors"
)

type OTPUsecase struct {
	OTPGenerator  ports.OTPGenerator
	OTPRepository ports.OTPRepository
	Cryptor       ports.Cryptor
}

func NewOTPUsecase(OTPGenerator ports.OTPGenerator, OTPRepository ports.OTPRepository, cryptor ports.Cryptor) *OTPUsecase {
	return &OTPUsecase{
		OTPGenerator:  OTPGenerator,
		OTPRepository: OTPRepository,
		Cryptor:       cryptor,
	}
}

func (o *OTPUsecase) SendOTP(emailDTO *dtos.Email) (*errors.HTTPError, *string) {
	OTP := o.OTPGenerator.Generate(6, false, false, false)

	now := time.Now()
	expirationTime := now.Add(time.Minute * 10)

	err, instanceId := o.OTPRepository.Create(OTP, expirationTime)
	if err != nil {
		return err, nil
	}

	if instanceId == nil {
		err := errors.NewRepositoryError(fmt.Errorf("ID from repository instance can not be empty"))
		return err, nil
	}

	details := &models.OTPDetails{
		ID:        *instanceId,
		Timestamp: now,
		Email:     emailDTO.Email,
		Success:   true,
		Message:   "OTP sent to user",
	}

	detailsMap := structs.Map(details)
	detailsJSON, parseErr := json.Marshal(detailsMap)
	if parseErr != nil {
		err := errors.NewValidationError(parseErr)
		return err, nil
	}

	err, encodedDetailsJSON := o.Cryptor.Encrypt(string(detailsJSON))
	if err != nil {
		return err, nil
	}

	// TODO: send email

	return nil, encodedDetailsJSON
}

func (o *OTPUsecase) VerifyOTP(OTPVerificationDTO dtos.OTPVerification) *errors.HTTPError {
	now := time.Now()

	err, decodedDetailsJSON := o.Cryptor.Decrypt(OTPVerificationDTO.EncryptedVerification)
	if err != nil {
		return err
	}

	var decodedDetails models.OTPDetails
	parseErr := json.Unmarshal([]byte(*decodedDetailsJSON), &decodedDetails)
	if parseErr != nil {
		err := errors.NewValidationError(parseErr)
		return err
	}

	emailFromDecodedDetails := decodedDetails.Email
	if OTPVerificationDTO.Email != emailFromDecodedDetails {
		err := errors.NewValidationError(fmt.Errorf("OTP was not sent to this email"))
		return err
	}

	err, OTPInstance := o.OTPRepository.Get(decodedDetails.ID)
	if err != nil {
		return err
	}

	if OTPInstance == nil {
		err := errors.NewRepositoryError(fmt.Errorf("OTP instance from repository can not be empty"))
		return err
	}

	if OTPInstance.Verified {
		err := errors.NewValidationError(fmt.Errorf("OTP is already verified"))
		return err
	}

	if OTPInstance.ExpirationTime.Before(now) {
		err := errors.NewValidationError(fmt.Errorf("OTP is expired"))

		repositoryErr := o.OTPRepository.Delete(decodedDetails.ID)
		if repositoryErr != nil {
			return repositoryErr
		}

		return err
	}

	if OTPInstance.OTP != OTPVerificationDTO.OTP {
		err := errors.NewValidationError(fmt.Errorf("OTP not matched"))
		return err
	}

	err = o.OTPRepository.Verify(decodedDetails.ID)
	if err != nil {
		return err
	}

	return nil
}
