package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/wazwki/skillsrock/docs"
	"github.com/wazwki/skillsrock/internal/app"
	"github.com/wazwki/skillsrock/internal/config"

	"github.com/joho/godotenv"
)

// @title Skillsrock API
// @version 1.0
// @description Skillsrock API documentation
// @termsOfService http://swagger.io/terms/

// @host localhost
// @BasePath /v1

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("[ERROR] Can't load environment: %s", err.Error())
	}
}

func main() {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("[ERROR] Can't load config: %s", err.Error())
	}

	app, err := app.New(cfg)
	if err != nil {
		log.Fatalf("[ERROR] Can't create app: %s", err.Error())
	}

	go func() {
		if err := app.Run(); err != nil {
			log.Fatalf("[ERROR] Can't run app: %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Print("Graceful shutdown start...")

	err = app.Stop()
	if err != nil {
		log.Fatalf("[ERROR] Can't gracefully close app: %s", err.Error())
	}
	log.Print("Graceful shutdown end...")
}
