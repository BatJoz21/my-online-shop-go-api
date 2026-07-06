package models

import (
	"github.com/BatJoz21/my-online-shop-go-api/database"
	"github.com/shopspring/decimal"
)

type Product struct {
	ID          int64           `json:"id"`
	CategoryID  int64           `json:"category_id"`
	Name        string          `json:"name"`
	Slug        string          `json:"slug"`
	Description string          `json:"description"`
	Price       decimal.Decimal `json:"price"`
	Stock       int64           `json:"stock"`
	Image       *string         `json:"image"`
	IsActive    bool            `json:"is_active"`
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
