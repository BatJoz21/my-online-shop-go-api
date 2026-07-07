package models

import (
	"github.com/BatJoz21/my-online-shop-go-api/database"
	"github.com/shopspring/decimal"
)

const (
	ProductPerPageLimit = 5
)

type Product struct {
	ID           int64           `json:"id"`
	CategoryID   int64           `json:"category_id"`
	Name         string          `json:"name"`
	Slug         string          `json:"slug"`
	Description  string          `json:"description"`
	Price        decimal.Decimal `json:"price"`
	Stock        int64           `json:"stock"`
	Image        *string         `json:"image"`
	IsActive     bool            `json:"is_active"`
	CategoryName string          `json:"category_name"`
}

func (p *Product) Save() error {
	query := `INSERT INTO products(category_id, name, slug, description, price, stock, image)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(p.CategoryID, p.Name, p.Slug, p.Description, p.Price, p.Stock, p.Image)
	if err != nil {
		return err
	}

	p.ID, err = result.LastInsertId()
	p.IsActive = true
	if err != nil {
		return err
	}

	return nil
}

func GetAllProducts(category string, offset int) ([]Product, error) {
	query := `SELECT
		products.id,
		products.category_id,
		products.name,
		products.slug,
		products.description,
		products.price,
		products.stock,
		products.image,
		products.is_active,
		categories.name AS category_name
	FROM products
	JOIN categories ON products.category_id = categories.id
	WHERE products.stock > 0 LIMIT ? OFFSET ?`

	rows, err := database.DB.Query(query, ProductPerPageLimit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		err = rows.Scan(&product.ID, &product.CategoryID, &product.Name, &product.Slug,
			&product.Description, &product.Price, &product.Stock, &product.Image,
			&product.IsActive, &product.CategoryName)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

func GetProduct(id int64) (*Product, error) {
	query := `SELECT
		products.id,
		products.category_id,
		products.name,
		products.slug,
		products.description,
		products.price,
		products.stock,
		products.image,
		products.is_active,
		categories.name AS category_name
	FROM products
	JOIN categories ON products.category_id = categories.id
	WHERE products.id = ?`
	row := database.DB.QueryRow(query, id)

	var product Product
	err := row.Scan(&product.ID, &product.CategoryID, &product.Name, &product.Slug,
		&product.Description, &product.Price, &product.Stock, &product.Image,
		&product.IsActive, &product.CategoryName)
	if err != nil {
		return nil, err
	}

	return &product, nil
}
