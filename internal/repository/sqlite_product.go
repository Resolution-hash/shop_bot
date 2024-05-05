package repository

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
	_ "github.com/mattn/go-sqlite3"
)

type SqliteProductRepo struct {
	db *sql.DB
}

func NewSqliteProductRepo(db *sql.DB) *SqliteProductRepo {
	return &SqliteProductRepo{db: db}
}

func (repo *SqliteProductRepo) DeleteProduct(id int64) error {
	_, err := prepareQuery("delete", "products", id).(squirrel.DeleteBuilder).RunWith(repo.db).Exec()
	if err != nil {
		return err
	}
	return nil
}

func (repo *SqliteProductRepo) UpdateProduct(product Product) error {
	_, err := prepareQuery("update", "products", product).(squirrel.UpdateBuilder).RunWith(repo.db).Exec()
	if err != nil {
		return err
	}
	return nil
}

func (repo *SqliteProductRepo) CreateProduct(product Product) error {

	_, err := prepareQuery("insert", "products", product).(squirrel.InsertBuilder).
		RunWith(repo.db).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

func (repo *SqliteProductRepo) GetAllProducts() ([]Product, error) {
	rows, err := prepareQuery("select", "products", nil).(squirrel.SelectBuilder).
		RunWith(repo.db).
		Query()
	if err != nil {
		return []Product{}, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		err = rows.Scan(&p.ID, &p.Name, &p.Type, &p.Description, &p.Price)
		if err != nil {
			return []Product{}, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (repo *SqliteProductRepo) GetProductsByType(productType string) ([]Product, error) {
	rows, err := prepareQuery("selectByType", "products", productType).(squirrel.SelectBuilder).
		RunWith(repo.db).
		Query()
	if err != nil {
		return []Product{}, err
	}
	var products []Product
	for rows.Next() {
		var p Product
		err = rows.Scan(&p.ID, &p.Name, &p.Type, &p.Description, &p.Price)
		if err != nil {
			return []Product{}, err
		}
		products = append(products, p)
	}
	return products, nil

}

func prepareQuery(operation string, table string, data interface{}) squirrel.Sqlizer {
	switch operation {
	case "insert":
		product := (data).(Product)
		insertMap := map[string]interface{}{
			"name":        product.Name,
			"type":        product.Type,
			"description": product.Description,
			"price":       product.Price,
		}
		return squirrel.Insert(table).SetMap(insertMap)
	case "select":
		return squirrel.Select("*").From(table)
	case "selectByType":
		return squirrel.Select("*").From(table).Where(squirrel.Eq{"type": data.(string)})
	case "update":
		product := (data).(Product)
		updateMap := map[string]interface{}{
			"name":        product.Name,
			"type":        product.Type,
			"description": product.Description,
			"price":       product.Price,
		}

		return squirrel.Update(table).SetMap(updateMap).Where(squirrel.Eq{"id": product.ID})
	case "delete":
		productID := (data).(int64)
		return squirrel.Delete(table).Where(squirrel.Eq{"id": productID})
	default:
		return nil
	}
}
