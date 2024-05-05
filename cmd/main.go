package main

import (
	"log"
	"telegram_bot/config"
	"telegram_bot/internal/app"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	app.StartBot(cfg)

}
