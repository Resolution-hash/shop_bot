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

// AddAdmin(int64) error
// 	RemoveAdmin(int64) error
// 	IsAdmin(int64) bool

func (repo *SqliteAdminRepo) AddAdmin(userID int64) error {
	_, err := prepareQueryAdmin("addAdmin", "admins", userID).(squirrel.InsertBuilder).
		RunWith(repo.Db).
		Exec()
	if err != nil {
		return err
	}
	return nil
}

func (repo *SqliteAdminRepo) IsAdmin(userID int64) (bool, error) {
	var count int
	err := prepareQueryAdmin("isAdmin", "admins", userID).(squirrel.SelectBuilder).
		RunWith(repo.Db).QueryRow().Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 1, nil
}

func prepareQueryAdmin(operation string, table string, data interface{}) squirrel.Sqlizer {
	switch operation {
	case "addAdmin":
		userID := data.(int64)
		return squirrel.Insert(table).Columns("admin_id").Values(userID)
	case "deleteAdmin":
		userID := data.(int64)
		return squirrel.Delete(table).Where(squirrel.Eq{"admin_id": userID})
	case "isAdmin":
		userID := data.(int64)
		return squirrel.Select("COUNT(*)").From(table).Where(squirrel.Eq{"admin_id": userID})
	}
	return nil
}
