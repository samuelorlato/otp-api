package main

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/samuelorlato/otp-api.git/internal/config"
	"github.com/samuelorlato/otp-api.git/internal/core/usecases"
	"github.com/samuelorlato/otp-api.git/internal/handlers"
	"github.com/samuelorlato/otp-api.git/internal/repositories"
	"github.com/samuelorlato/otp-api.git/pkg/cryptography"
	"github.com/samuelorlato/otp-api.git/pkg/email"
	"github.com/samuelorlato/otp-api.git/pkg/otp"
)

func main() {
	ctx := context.Background()

	// Uncomment to run locally
	// err := godotenv.Load()
	// if err != nil {
	// 	panic(err)
	// }

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

	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	from := os.Getenv("FROM_EMAIL")
	password := os.Getenv("FROM_EMAIL_PASSWORD")
	emailSender := email.NewEmailSender(host, port, from, password)

	OTPUsecase := usecases.NewOTPUsecase(OTPGenerator, OTPRepository, cryptor, emailSender)

	HTTPHandler := handlers.NewHTTPHandler(engine, OTPUsecase, errorHandler)
	HTTPHandler.SetRoutes()

	engine.Run()
}
