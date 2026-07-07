package models

import (
	"github.com/BatJoz21/my-online-shop-go-api/database"
	"github.com/shopspring/decimal"
)

type ProductVariant struct {
	ID            int64           `json:"id"`
	ProductID     int64           `json:"product_id"`
	Name          string          `json:"name"`
	Sku           string          `json:"sku"`
	PriceModifier decimal.Decimal `json:"price_modifier"`
}

func (pV *ProductVariant) Save() error {
	query := `INSERT INTO product_variants(product_id, name, sku, price_modifier)
		VALUES (?, ?, ?, ?)`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(pV.ProductID, pV.Name, pV.Sku, pV.PriceModifier)
	if err != nil {
		return err
	}

	pV.ID, err = result.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}
