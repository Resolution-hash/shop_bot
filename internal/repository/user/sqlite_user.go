package repository

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
)

type SqliteUserRepo struct {
	Db *sql.DB
}

func NewSqliteUserRepo(db *sql.DB) *SqliteUserRepo {
	return &SqliteUserRepo{
		Db: db,
	}
}

func (repo *SqliteUserRepo) AddUser(user User) error {
	_, err := prepareQueryUser("addUser", "users", user).(squirrel.InsertBuilder).
		RunWith(repo.Db).
		Exec()
	if err != nil {
		return err
	}
	return nil
}

func prepareQueryUser(operation string, table string, data interface{}) squirrel.Sqlizer {
	switch operation {
	case "addUser":
		user := data.(User)
		insertMap := map[string]interface{}{
			"id":         user.UserID,
			"first_name": user.First_name,
			"user_name":  user.User_name,
		}

		return squirrel.Insert(table).SetMap(insertMap).Suffix("ON CONFLICT(id) DO REPLACE")
	}
	return nil
}
