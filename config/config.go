package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var cfgPath = "C:/Users/Leoni/Documents/projects/shop_bot/config/.env"

type Config struct {
	TelegramAPIToken string
	DbUrl            string
	ImagesUrl        string
}

func LoadConfig() (*Config, error) {

	err := godotenv.Load(cfgPath)

	if err != nil {
		log.Fatal("Error loading .env file, ", err)
	}

	config := Config{
		TelegramAPIToken: os.Getenv("BOT_TOKEN"),
		DbUrl:            os.Getenv("DB_URL"),
		ImagesUrl:        os.Getenv("IMAGES_URL"),
	}

	return &config, nil
}
