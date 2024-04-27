package main

import (
	"log"
	"telegram_bot/config"
	"telegram_bot/internal/app"
)

func main() {
	config, err := config.LoadConfig("../config/config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	app.StartBot(config)

}
