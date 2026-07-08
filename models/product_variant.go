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

func GetAllVariantOfAProduct(product_id int64) (*[]ProductVariant, error) {
	query := `SELECT
		id,
		product_id,
		name,
		sku,
		price_modifier,
		stock
	FROM product_variants
	WHERE product_id = ?`
	rows, err := database.DB.Query(query, product_id)
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

func GetVariant(id int64) (*ProductVariant, error) {
	query := `SELECT
		id,
		product_id,
		name,
		sku,
		price_modifier,
		stock
	FROM product_variants
	WHERE id = ?`
	row := database.DB.QueryRow(query, id)

	var pv ProductVariant
	err := row.Scan(&pv.ID, &pv.ProductID, &pv.Name, &pv.Sku, &pv.PriceModifier, &pv.Stock)
	if err != nil {
		return nil, err
	}

	return &pv, nil
}

func (pV *ProductVariant) Update() error {
	query := `UPDATE product_variants SET
		product_id = ?,
		name = ?,
		sku = ?,
		price_modifier = ?,
		stock = ?
	WHERE id = ?`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(pV.ProductID, pV.Name, pV.Sku, pV.PriceModifier, pV.Stock, pV.ID)

	return err
}

func (pv *ProductVariant) UpdateStock() error {
	query := `UPDATE product_variants SET stock = ? WHERE id = ?`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(pv.Stock, pv.ID)

	return err
}

func DeleteAllVariantOfAProduct(product_id int64) error {
	query := `DELETE FROM product_variants WHERE product_id = ?`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return nil
	}

	_, err = stmt.Exec(product_id)

	return err
}

func (pv *ProductVariant) Delete() error {
	query := `DELETE FROM product_variants WHERE id = ?`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return nil
	}

	_, err = stmt.Exec(pv.ID)

	return err
}
