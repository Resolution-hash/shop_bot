package repository

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
)

type PostgersProductRepo struct {
	db *sql.DB
}

func NewSqliteProductRepo(db *sql.DB) *PostgersProductRepo {
	return &PostgersProductRepo{db: db}
}

func (repo *PostgersProductRepo) DeleteProduct(id int64) error {
	_, err := PrepareQueryProduct("delete", "products", id).(squirrel.DeleteBuilder).RunWith(repo.db).Exec()
	if err != nil {
		return err
	}
	return nil
}

func (repo *PostgersProductRepo) UpdateProduct(product Product) error {
	_, err := PrepareQueryProduct("update", "products", product).(squirrel.UpdateBuilder).RunWith(repo.db).Exec()
	if err != nil {
		return err
	}
	return nil
}

func (repo *PostgersProductRepo) CreateProduct(product Product) error {
	_, err := PrepareQueryProduct("insert", "products", product).(squirrel.InsertBuilder).
		RunWith(repo.db).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

// func (repo *PostgersProductRepo) CreateTestProduct(product Product) error {
// 	_, err := PrepareQueryProduct("insert", "test_products", product).(squirrel.InsertBuilder).
// 		RunWith(repo.db).
// 		Exec()
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (repo *PostgersProductRepo) GetAllProducts() ([]Product, error) {
	rows, err := PrepareQueryProduct("select", "products", nil).(squirrel.SelectBuilder).
		RunWith(repo.db).
		Query()
	if err != nil {
		return []Product{}, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		err = rows.Scan(&p.ID, &p.Name, &p.Type, &p.Description, &p.Price, &p.Image)
		if err != nil {
			return []Product{}, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (repo *PostgersProductRepo) GetProductsByType(productType string) ([]Product, error) {
	rows, err := PrepareQueryProduct("selectByType", "products", productType).(squirrel.SelectBuilder).
		RunWith(repo.db).
		Query()
	if err != nil {
		return []Product{}, err
	}
	var products []Product
	for rows.Next() {
		var p Product
		err = rows.Scan(&p.ID, &p.Name, &p.Type, &p.Description, &p.Price, &p.Image)
		if err != nil {
			return []Product{}, err
		}
		products = append(products, p)
	}
	return products, nil

}

func PrepareQueryProduct(operation string, table string, data interface{}) squirrel.Sqlizer {
	switch operation {
	case "insert":
		product := (data).(Product)
		insertMap := map[string]interface{}{
			"name":        product.Name,
			"type":        product.Type,
			"description": product.Description,
			"price":       product.Price,
			"image":       product.Image,
		}
		return squirrel.Insert(table).SetMap(insertMap).PlaceholderFormat(squirrel.Dollar)
	case "select":
		return squirrel.Select("*").From(table).PlaceholderFormat(squirrel.Dollar)
	case "selectByType":
		return squirrel.Select("*").From(table).Where(squirrel.Eq{"type": data.(string)}).PlaceholderFormat(squirrel.Dollar)
	case "selectByIDs":
		return squirrel.Select("*").From(table).Where(squirrel.Eq{"id": data.([]int64)}).PlaceholderFormat(squirrel.Dollar)
	case "update":
		product := (data).(Product)
		updateMap := map[string]interface{}{
			"name":        product.Name,
			"type":        product.Type,
			"description": product.Description,
			"price":       product.Price,
			"image":       product.Image,
		}

		return squirrel.Update(table).SetMap(updateMap).Where(squirrel.Eq{"id": product.ID}).PlaceholderFormat(squirrel.Dollar)
	case "delete":
		productID := (data).(int64)
		return squirrel.Delete(table).Where(squirrel.Eq{"id": productID}).PlaceholderFormat(squirrel.Dollar)
	default:
		return nil
	}
}
