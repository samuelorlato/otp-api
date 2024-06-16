package repositories

import (
	"context"
	"encoding/json"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"github.com/samuelorlato/otp-api.git/internal/config"
	"github.com/samuelorlato/otp-api.git/internal/core/models"
	"github.com/samuelorlato/otp-api.git/internal/core/ports"
	"github.com/samuelorlato/otp-api.git/pkg/errors"
)

type FirestoreRepository struct {
	firestoreClient *firestore.Client
}

func NewFirestoreRepository(firestoreClient *firestore.Client) ports.OTPRepository {
	return &FirestoreRepository{
		firestoreClient: firestoreClient,
	}
}

func (f *FirestoreRepository) Create(OTP string, expirationTime time.Time) (*errors.HTTPError, *string) {
	instanceId := uuid.NewString()

	_, err := f.firestoreClient.Collection(config.FirestoreCollectionName).Doc(instanceId).Set(context.Background(), map[string]interface{}{
		"OTP":            OTP,
		"expirationTime": expirationTime,
		"verified":       false,
	})
	if err != nil {
		err := errors.NewRepositoryError(err)
		return err, nil
	}

	return nil, &instanceId
}

func (f *FirestoreRepository) Get(ID string) (*errors.HTTPError, *models.OTPInstance) {
	documentSnapshot, err := f.firestoreClient.Collection(config.FirestoreCollectionName).Doc(ID).Get(context.Background())
	if err != nil {
		err := errors.NewRepositoryError(err)
		return err, nil
	}

	OTPInstanceMap := documentSnapshot.Data()
	OTPInstanceBytes, parseErr := json.Marshal(OTPInstanceMap)
	if parseErr != nil {
		err := errors.NewValidationError(parseErr)
		return err, nil
	}

	var OTPInstance models.OTPInstance
	parseErr = json.Unmarshal(OTPInstanceBytes, &OTPInstance)
	if parseErr != nil {
		err := errors.NewValidationError(parseErr)
		return err, nil
	}

	return nil, &OTPInstance
}

func (f *FirestoreRepository) Delete(ID string) *errors.HTTPError {
	_, err := f.firestoreClient.Collection(config.FirestoreCollectionName).Doc(ID).Delete(context.Background())
	if err != nil {
		err := errors.NewRepositoryError(err)
		return err
	}

	return nil
}

func (f *FirestoreRepository) Verify(ID string) *errors.HTTPError {
	_, err := f.firestoreClient.Collection(config.FirestoreCollectionName).Doc(ID).Set(context.Background(), map[string]interface{}{
		"verified": true,
	}, firestore.MergeAll)
	if err != nil {
		err := errors.NewRepositoryError(err)
		return err
	}

	deleteErr := f.Delete(ID)
	if deleteErr != nil {
		return deleteErr
	}

	return nil
}
