package logic

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Resolution-hash/shop_bot/internal/services"
	db "github.com/Resolution-hash/shop_bot/repository/db"
	product "github.com/Resolution-hash/shop_bot/repository/product"
)

func ParseProduct(data string) (product.Product, error) {
	lines := strings.Split(data, "\n")
	if len(lines) != 4 {
		return product.Product{}, errors.New("Ошибка заполнения данных")
	}

	price, err := strconv.ParseFloat(strings.TrimSpace(lines[3]), 64)
	if err != nil {
		return product.Product{}, errors.New("Ошибка конвертации цены")
	}

	product := product.Product{
		Name:        strings.TrimSpace(lines[0]),
		Type:        strings.TrimSpace(lines[1]),
		Description: strings.TrimSpace(lines[2]),
		Price:       price,
	}

	return product, nil
}

func GetTestText(p product.Product) string {
	return fmt.Sprintf("%s\n\n%s\n\nЦена: %0.f рублей\n", p.Name, p.Description, p.Price)
}

func CreateProduct(product product.Product) error {
	db, err := db.SetupDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	service := InitProductService(db)

	err = service.CreateProduct(product)
	if err != nil {
		return err
	}
	return nil
}

func InitProductService(db *sql.DB) *services.ProductService {
	repo := product.NewSqliteProductRepo(db)
	service := services.NewProductService(repo)
	return service
}
