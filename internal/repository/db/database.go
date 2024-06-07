package repository

import (
	"database/sql"
	"fmt"

	"github.com/Resolution-hash/shop_bot/config"
)

func SetupDatabase() (*sql.DB, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	connect := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Db_name)

	db, err := sql.Open("postgres", connect)
	if err != nil {
		return nil, err
	}
	return db, nil
}
