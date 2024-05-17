package repository

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
)

type SqliteCartRepo struct {
	Db *sql.DB
}

func NewSqliteCartRepo(db *sql.DB) *SqliteCartRepo {
	return &SqliteCartRepo{
		Db: db,
	}
}
func (repo *SqliteCartRepo) AddItem(item CartItem) (int, error) {

	tx, err := repo.Db.Begin()
	if err != nil {
		return 0, err
	}
	_, err = prepareQueryCart("addItem", "cart", item).(squirrel.InsertBuilder).
		RunWith(repo.Db).
		Exec()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var total int
	err = prepareQueryCart("countByUserIDAndProductID", "cart", item).(squirrel.SelectBuilder).
		RunWith(repo.Db).
		QueryRow().Scan(&total)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return total, nil
}

// func RemoveItem(int64) error {

// }

// func GetAllItems() ([]CartItem, error){

// }

func prepareQueryCart(operation string, table string, data interface{}) squirrel.Sqlizer {
	switch operation {
	case "addItem":
		cartItem := (data).(CartItem)
		insertMap := map[string]interface{}{
			"product_id": cartItem.ProductID,
			"user_id":    cartItem.UserID,
			"quantity":   cartItem.Quantity,
		}
		return squirrel.Insert(table).SetMap(insertMap)
	case "selectByID":
		return squirrel.Select("*").From(table).Where(squirrel.Eq{"user_id": data.(string)})
	case "increment":
		return squirrel.Update(table).Set("quantity", squirrel.Expr("quantity + 1")).Where(squirrel.Eq{"product_id": data.(int64)})
	case "decrement":
		decrease := squirrel.Update(table).Set("quantity", squirrel.Expr("quantity - 1")).Where(squirrel.Eq{"product_id": data.(int64)}).Where("quantity > 1")
		delete := squirrel.Delete(table).Where(squirrel.Eq{"product_id": data.(int64)}).Where("quantity = 1")
		return squirrel.Case().When(squirrel.Expr("quantity > 1"), decrease).Else(delete)
	case "countByUserIDAndProductID":
		cartItem := data.(CartItem)
		return squirrel.Select("SUM(quantity) as total_quantity").From(table).Where(squirrel.Eq{"user_id": cartItem.UserID, "product_id": cartItem.ProductID})
	case "delete":
		return squirrel.Delete(table).Where(squirrel.Eq{"id": data.(int64)})
	default:
		return nil
	}
}
