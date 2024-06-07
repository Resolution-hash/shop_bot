package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

var cfgPath = "C:/Users/Leoni/Documents/projects/shop_bot/config/.env"

type Config struct {
	TelegramAPIToken string
	DbUrl            string
	ImagesUrl        string
	Host             string
	Port             string
	User             string
	Password         string
	Db_name          string
	Backet           string
	Endpoint         string
	AccessKey        string
	SecretKey        string
}

func LoadConfig() (*Config, error) {

	err := godotenv.Load(cfgPath)

	if err != nil {
		return nil, errors.New("error loading .env file: " + err.Error())
	}

	config := Config{
		TelegramAPIToken: os.Getenv("BOT_TOKEN"),
		DbUrl:            os.Getenv("DB_URL"),
		ImagesUrl:        os.Getenv("IMAGES_URL"),
		Host:             os.Getenv("HOST"),
		Port:             os.Getenv("PORT"),
		User:             os.Getenv("USER"),
		Password:         os.Getenv("PASSWORD"),
		Db_name:          os.Getenv("DB_NAME"),
		Backet:           os.Getenv("BACKET"),
		Endpoint:         os.Getenv("ENDPOINT"),
		AccessKey:        os.Getenv("ACCESS_KEY"),
		SecretKey:        os.Getenv("SECRET_KEY"),
	}

	return &config, nil
}
