package main

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/samuelorlato/otp-api.git/internal/config"
	"github.com/samuelorlato/otp-api.git/internal/core/usecases"
	"github.com/samuelorlato/otp-api.git/internal/handlers"
	"github.com/samuelorlato/otp-api.git/internal/repositories"
	"github.com/samuelorlato/otp-api.git/pkg/cryptography"
	"github.com/samuelorlato/otp-api.git/pkg/otp"
)

func main() {
	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	engine := gin.Default()

	errorHandler := handlers.NewErrorHandler()

	OTPGenerator := otp.NewOTPGenerator()

	credentialsJSONString := os.Getenv("FIREBASE_CREDENTIALS")
	app, err := config.InitFirebaseApp(ctx, credentialsJSONString)
	if err != nil {
		panic(err)
	}

	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		panic(err)
	}
	OTPRepository := repositories.NewFirestoreRepository(firestoreClient)

	key := os.Getenv("KEY")
	iv := os.Getenv("IV")
	cryptor := cryptography.NewCryptor(key, iv)

	OTPUsecase := usecases.NewOTPUsecase(OTPGenerator, OTPRepository, cryptor)

	HTTPHandler := handlers.NewHTTPHandler(engine, OTPUsecase, errorHandler)
	HTTPHandler.SetRoutes()

	engine.Run()
}
