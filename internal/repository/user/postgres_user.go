package repository

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/gookit/color"
)

type PostgresUserRepo struct {
	Db *sql.DB
}

func NewPostgresUserRepo(db *sql.DB) *PostgresUserRepo {
	return &PostgresUserRepo{
		Db: db,
	}
}

func (repo *PostgresUserRepo) AddUser(user User) (bool, error) {
	var isAdmin int
	err := prepareQueryUser("isAdmin", "admins", user.UserID).(squirrel.SelectBuilder).
		RunWith(repo.Db).QueryRow().Scan(&isAdmin)
	if err != nil {
		return isAdmin == 1, err
	}
	color.Redln("User ", user.UserID, "isAdmin ", isAdmin)
	usr := User{
		UserID:     user.UserID,
		First_name: user.First_name,
		Last_name:  user.Last_name,
		User_name:  user.User_name,
		IsAdmin:    isAdmin,
	}
	_, err = prepareQueryUser("addUser", "users", usr).(squirrel.InsertBuilder).
		RunWith(repo.Db).
		Exec()
	if err != nil {
		return false, err
	}
	return isAdmin == 1, nil
}

func prepareQueryUser(operation string, table string, data interface{}) squirrel.Sqlizer {
	switch operation {
	case "addUser":
		user := data.(User)
		insertMap := map[string]interface{}{
			"id":         user.UserID,
			"first_name": user.First_name,
			"user_name":  user.User_name,
			"is_admin":   user.IsAdmin,
		}
		return squirrel.Insert(table).SetMap(insertMap).Suffix("ON CONFLICT(id) DO NOTHING").PlaceholderFormat(squirrel.Dollar)
	case "isAdmin":
		userID := data.(int)
		return squirrel.Select("COUNT(*)").From(table).Where(squirrel.Eq{"admin_id": userID}).PlaceholderFormat(squirrel.Dollar)
	}
	return nil
}
