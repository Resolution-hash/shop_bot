package repository

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
)

type SqliteAdminRepo struct {
	Db *sql.DB
}

func NewSqliteAdminRepo(db *sql.DB) *SqliteAdminRepo {
	return &SqliteAdminRepo{
		Db: db,
	}
}

func (repo *SqliteAdminRepo) AddAdmin(userID int64) error {
	_, err := PrepareQueryAdmin("addAdmin", "admins", userID).(squirrel.InsertBuilder).
		RunWith(repo.Db).
		Exec()
	if err != nil {
		return err
	}
	return nil
}

func (repo *SqliteAdminRepo) IsAdmin(userID int64) (bool, error) {
	var count int
	err := PrepareQueryAdmin("isAdmin", "admins", userID).(squirrel.SelectBuilder).
		RunWith(repo.Db).QueryRow().Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 1, nil
}

func PrepareQueryAdmin(operation string, table string, data interface{}) squirrel.Sqlizer {
	switch operation {
	case "addAdmin":
		userID := data.(int)
		return squirrel.Insert(table).Columns("admin_id").Values(userID).PlaceholderFormat(squirrel.Dollar)
	case "deleteAdmin":
		userID := data.(int)
		return squirrel.Delete(table).Where(squirrel.Eq{"admin_id": userID}).PlaceholderFormat(squirrel.Dollar)
	case "isAdmin":
		userID := data.(int)
		return squirrel.Select("COUNT(*)").From(table).Where(squirrel.Eq{"admin_id": userID}).PlaceholderFormat(squirrel.Dollar)
	}
	return nil
}
