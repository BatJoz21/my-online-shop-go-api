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
	Stock         int64           `json:"stock"`
}

func (pV *ProductVariant) Save() error {
	query := `INSERT INTO product_variants(product_id, name, sku, price_modifier, stock)
		VALUES (?, ?, ?, ?, ?)`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(pV.ProductID, pV.Name, pV.Sku, pV.PriceModifier, pV.Stock)
	if err != nil {
		return err
	}

	pV.ID, err = result.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

func GetAllVariantOfAProduct(id int64) (*[]ProductVariant, error) {
	query := `SELECT
		id,
		product_id,
		name,
		sku,
		price_modifier,
		stock
	FROM product_variants
	WHERE product_id = ?`
	rows, err := database.DB.Query(query, id)
	if err != nil {
		return nil, err
	}

	var variants []ProductVariant
	for rows.Next() {
		var variant ProductVariant
		err = rows.Scan(&variant.ID, &variant.ProductID, &variant.Name,
			&variant.Sku, &variant.PriceModifier, &variant.Stock)
		if err != nil {
			return nil, err
		}

		variants = append(variants, variant)
	}

	return &variants, nil
}
