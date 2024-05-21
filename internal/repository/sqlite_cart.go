package repository

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/gookit/color"
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
	err = prepareQueryCart("countByID", "cart", item).(squirrel.SelectBuilder).
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

func (repo *SqliteCartRepo) Increment(item CartItem) (int, error) {
	tx, err := repo.Db.Begin()
	if err != nil {
		return 0, err
	}
	color.Redln("productID:", item.ProductID, "UserID:", item.UserID)

	_, err = prepareQueryCart("increment", "cart", item).(squirrel.UpdateBuilder).
		RunWith(repo.Db).
		Exec()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var total int
	err = prepareQueryCart("countByID", "cart", item).(squirrel.SelectBuilder).
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

func (repo *SqliteCartRepo) Decrement(item CartItem) (int, error) {
	tx, err := repo.Db.Begin()
	if err != nil {
		return 0, err
	}

	var total int
	err = prepareQueryCart("countByID", "cart", item).(squirrel.SelectBuilder).
		RunWith(repo.Db).
		QueryRow().Scan(&total)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if total > 1 {
		_, err = prepareQueryCart("decrement", "cart", item).(squirrel.UpdateBuilder).
			RunWith(repo.Db).
			Exec()
		if err != nil {
			return 0, err
		}
		total -= 1
	} else if total == 1 {
		_, err = prepareQueryCart("delete", "cart", item).(squirrel.DeleteBuilder).
			RunWith(repo.Db).
			Exec()
		if err != nil {
			return 0, err
		}
		total = 0
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return total, nil
}

// RemoveItem(int64) error
// GetItemsByID(int64) ([]CartItem, error)

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
		cartItem := data.(CartItem)
		color.Redln("productID:", cartItem.ProductID, "UserID:", cartItem.UserID)
		return squirrel.Update(table).Set("quantity", squirrel.Expr("quantity + 1")).Where(squirrel.Eq{"product_id": cartItem.ProductID, "user_id": cartItem.UserID})
	case "decrement":
		cartItem := data.(CartItem)
		return squirrel.Update(table).Set("quantity", squirrel.Expr("quantity - 1")).Where(squirrel.Eq{"user_id": cartItem.UserID, "product_id": cartItem.ProductID}).Where("quantity > 1")
	case "countByID":
		cartItem := data.(CartItem)
		return squirrel.Select("SUM(quantity) as total_quantity").From(table).Where(squirrel.Eq{"user_id": cartItem.UserID, "product_id": cartItem.ProductID})
	case "delete":
		cartItem := data.(CartItem)
		return squirrel.Delete(table).Where(squirrel.Eq{"user_id": cartItem.UserID, "product_id": cartItem.ProductID})
	default:
		return nil
	}
}
