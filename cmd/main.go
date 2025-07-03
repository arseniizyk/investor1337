package main

import (
	"log"

	"github.com/arseniizyk/investor1337/internal/app"
	"github.com/arseniizyk/investor1337/internal/config"
	"go.uber.org/zap"
)

func main() {
	l, _ := zap.NewDevelopment()
	defer l.Sync()

	l.Info("Initializing config")
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	l.Info("Config initialized")

	a := app.New(cfg, l)

	l.Info("Running app")
	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}
