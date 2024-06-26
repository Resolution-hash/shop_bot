package repository

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	product "github.com/Resolution-hash/shop_bot/repository/product"

	"github.com/gookit/color"
)

type PostgresCartRepo struct {
	Db *sql.DB
}

func NewSqliteCartRepo(db *sql.DB) *PostgresCartRepo {
	return &PostgresCartRepo{
		Db: db,
	}
}

func (repo *PostgresCartRepo) GetQuantityByItemID(item CartItem) (int, error) {
	var total int
	err := prepareQueryProductCart("quantityByID", "cart", item).(squirrel.SelectBuilder).
		RunWith(repo.Db).
		QueryRow().Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (repo *PostgresCartRepo) AddItem(item CartItem) (int, error) {

	tx, err := repo.Db.Begin()
	if err != nil {
		return 0, err
	}
	_, err = prepareQueryProductCart("addItem", "cart", item).(squirrel.InsertBuilder).
		RunWith(repo.Db).
		Exec()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var total int
	err = prepareQueryProductCart("quantityByID", "cart", item).(squirrel.SelectBuilder).
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

func (repo *PostgresCartRepo) Increment(item CartItem) (int, error) {
	tx, err := repo.Db.Begin()
	if err != nil {
		return 0, err
	}
	color.Redln("productID:", item.ProductID, "UserID:", item.UserID)

	_, err = prepareQueryProductCart("increment", "cart", item).(squirrel.UpdateBuilder).
		RunWith(repo.Db).
		Exec()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var total int
	err = prepareQueryProductCart("quantityByID", "cart", item).(squirrel.SelectBuilder).
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

func (repo *PostgresCartRepo) Decrement(item CartItem) (int, error) {
	var total int
	err := prepareQueryProductCart("quantityByID", "cart", item).(squirrel.SelectBuilder).
		RunWith(repo.Db).
		QueryRow().Scan(&total)
	if err != nil {
		return 0, err
	}

	if total > 1 {
		_, err = prepareQueryProductCart("decrement", "cart", item).(squirrel.UpdateBuilder).
			RunWith(repo.Db).
			Exec()
		if err != nil {
			return 0, err
		}
		total -= 1
	} else if total == 1 {
		_, err = prepareQueryProductCart("delete", "cart", item).(squirrel.DeleteBuilder).
			RunWith(repo.Db).
			Exec()
		if err != nil {
			return 0, err
		}
		total = 0
	}

	return total, nil
}

func (repo *PostgresCartRepo) GetItemsByUserID(userID int64) ([]*CartProduct, error) {

	items := make([]*CartItem, 0)

	rows, err := prepareQueryProductCart("selectByID", "cart", userID).(squirrel.SelectBuilder).
		RunWith(repo.Db).
		Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := new(CartItem)
		if err := rows.Scan(&item.ProductID, &item.Quantity); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	productIDs := make([]int64, len(items))
	for _, item := range items {
		productIDs = append(productIDs, item.ProductID)
	}

	rows, err = product.PrepareQueryProduct("selectByIDs", "products", productIDs).(squirrel.SelectBuilder).
		RunWith(repo.Db).
		Query()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make(map[int64]*product.Product, len(productIDs))
	for rows.Next() {
		product := new(product.Product)
		if err := rows.Scan(&product.ID, &product.Name, &product.Type, &product.Description, &product.Price, &product.Image); err != nil {
			return nil, err
		}
		products[product.ID] = product
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	cartProducts := make([]*CartProduct, len(items))
	for i, item := range items {
		product, ok := products[item.ProductID]
		if !ok {
			return nil, fmt.Errorf("product not found: %d", item.ProductID)
		}
		cartProducts[i] = &CartProduct{
			ProductID:   int(product.ID),
			Name:        product.Name,
			Type:        product.Type,
			Description: product.Description,
			Price:       product.Price,
			Quantity:    item.Quantity,
			Image:       product.Image,
		}

	}

	return cartProducts, nil
}

func (repo *PostgresCartRepo) DeleteItem(item CartItem) error {
	_, err := prepareQueryProductCart("delete", "cart", item).(squirrel.DeleteBuilder).
		RunWith(repo.Db).
		Exec()
	if err != nil {
		return err
	}
	return nil
}

func prepareQueryProductCart(operation string, table string, data interface{}) squirrel.Sqlizer {
	switch operation {
	case "addItem":
		cartItem := (data).(CartItem)
		insertMap := map[string]interface{}{
			"product_id": cartItem.ProductID,
			"user_id":    cartItem.UserID,
			"quantity":   cartItem.Quantity,
		}
		return squirrel.Insert(table).SetMap(insertMap).PlaceholderFormat(squirrel.Dollar)
	case "selectByID":
		return squirrel.Select("product_id, quantity").From(table).Where(squirrel.Eq{"user_id": data.(int64)}).PlaceholderFormat(squirrel.Dollar)
	case "increment":
		cartItem := data.(CartItem)
		color.Redln("productID:", cartItem.ProductID, "UserID:", cartItem.UserID)
		return squirrel.Update(table).Set("quantity", squirrel.Expr("quantity + 1")).Where(squirrel.Eq{"product_id": cartItem.ProductID, "user_id": cartItem.UserID}).PlaceholderFormat(squirrel.Dollar)
	case "decrement":
		cartItem := data.(CartItem)
		return squirrel.Update(table).Set("quantity", squirrel.Expr("quantity - 1")).Where(squirrel.Eq{"user_id": cartItem.UserID, "product_id": cartItem.ProductID}).Where("quantity > 1").PlaceholderFormat(squirrel.Dollar)
	case "quantityByID":
		cartItem := data.(CartItem)
		return squirrel.Select("COALESCE(SUM(quantity), 0) as total_quantity").From(table).Where(squirrel.Eq{"user_id": cartItem.UserID, "product_id": cartItem.ProductID}).PlaceholderFormat(squirrel.Dollar)
	case "isTableEmpty":
		cartItem := data.(CartItem)
		return squirrel.Select("COUNT(*)").From(table).Where(squirrel.Eq{"user_id": cartItem.UserID}).PlaceholderFormat(squirrel.Dollar)
	case "delete":
		cartItem := data.(CartItem)
		return squirrel.Delete(table).Where(squirrel.Eq{"user_id": cartItem.UserID, "product_id": cartItem.ProductID}).PlaceholderFormat(squirrel.Dollar)
	default:
		return nil
	}
}
