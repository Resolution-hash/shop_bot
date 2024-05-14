package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var cfgPath = "C:/Users/Leoni/OneDrive/Документы/projects/shop_bot/config/.env"

type Config struct {
	TelegramAPIToken string
	DbUrl            string
}

func LoadConfig() (*Config, error) {

	err := godotenv.Load(cfgPath)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := Config{
		TelegramAPIToken: os.Getenv("BOT_TOKEN"),
		DbUrl:            os.Getenv("DB_URL"),
	}

	return &config, nil
}
