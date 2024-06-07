package logic

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	// db "github.com/Resolution-hash/shop_bot/internal/repository/db"
	product "github.com/Resolution-hash/shop_bot/internal/repository/product"
)

// func AddTestProduct(productItem product.Product) error {
// 	db, err := db.SetupDatabase()
// 	if err != nil {
// 		return err
// 	}
// 	defer db.Close()

// }

// func showTestProduct

func ParseProduct(data string) (product.Product, error) {
	lines := strings.Split(data, "\n")
	if len(lines) != 5 {
		return product.Product{}, errors.New("неверный формат. Строк должно быть 5")
	}

	price, err := strconv.ParseFloat(strings.TrimSpace(lines[3]), 16)
	if err != nil {
		return product.Product{}, errors.New("ошибка конвертации цены")
	}

	product := product.Product{
		Name:        strings.TrimSpace(lines[0]),
		Type:        strings.TrimSpace(lines[1]),
		Description: strings.TrimSpace(lines[2]),
		Price:       price,
		Image:       strings.TrimSpace(lines[4]),
	}

	return product, nil
}


func GetTestText(p product.Product) string {
	return fmt.Sprintf("%s\n\n%s\n\nЦена: %0.f рублей\n", p.Name, p.Description, p.Price)
}


