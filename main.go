package main

import (
	"github.com/RevinB/mira_server/config"
	"github.com/RevinB/mira_server/data"
	"github.com/RevinB/mira_server/handler"
	"github.com/RevinB/mira_server/router"
	"github.com/RevinB/mira_server/utils"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/getsentry/sentry-go"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Println("Welcome! Starting server...")

	// Log into sentry
	err := sentry.Init(sentry.ClientOptions{Dsn: os.Getenv("SENTRY_DSN")})
	if err != nil {
		panic("error connecting to sentry: " + err.Error())
	}

	cfg := config.Config{
		AppUrl:       os.Getenv("APP_ADDR"),
		FinalUrlBase: os.Getenv("FINAL_URL_BASE"),
		S3BucketName: os.Getenv("AWS_S3_BUCKET_NAME"),
		JWTSecret:    utils.GetenvByteArray("JWT_SECRET"),
	}

	// new db conn
	db, err := data.NewStore()
	if err != nil {
		panic("database connection failed: " + err.Error())
	}
	log.Println("Database connection successful")

	err = data.Migrate(db.Client)
	if err != nil {
		panic("database migration failed: " + err.Error())
	}
	log.Println("Database migration successful")

	awsSession := session.Must(session.NewSession())

	iRouter := router.NewRouter()

	iHandler := handler.NewHandler(db, cfg, awsSession)
	iHandler.ImplHandler(iRouter)

	go func() {
		if err := iRouter.Listen(cfg.AppUrl); err != nil && err != http.ErrServerClosed {
			sentry.CaptureException(err)
			panic("failed to initialize server: " + err.Error())
		}
	}()

	log.Println("API started. Listening for requests.")

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	if err := iRouter.Shutdown(); err != nil {
		sentry.CaptureException(err)
		panic("failed to gracefully shutdown fiber: " + err.Error())
	}
}
