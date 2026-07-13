package models

import (
	"github.com/BatJoz21/my-online-shop-go-api/database"
	"github.com/shopspring/decimal"
)

type OrderItem struct {
	ID            int64           `json:"id"`
	OrderID       int64           `json:"order_id"`
	ProductID     int64           `json:"product_id"`
	VariantID     int64           `json:"variant_id"`
	ProductName   string          `json:"product_name_snapshot"`
	Quantity      int             `json:"quantity"`
	PriceSnapshot decimal.Decimal `json:"price_snapshot"`
	Subtotal      decimal.Decimal `json:"subtotal"`
}

func (o *OrderItem) Save() error {
	query := `INSERT INTO order_items(
		order_id,
		product_id,
		variant_id,
		product_name_snapshot,
		quantity,
		price_snapshot,
		subtotal)
	VALUES (?, ?, ?, ?, ?, ?, ?)`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(o.OrderID, o.ProductID, o.VariantID, o.ProductName, o.Quantity,
		o.PriceSnapshot, o.Subtotal)
	if err != nil {
		return err
	}

	o.ID, err = result.LastInsertId()
	return err
}

func GetAllItemFromOrder(orderID int64) (*[]OrderItem, error) {
	query := `SELECT
		id,
		order_id,
		product_id,
		variant_id,
		product_name_snapshot,
		quantity,
		price_snapshot,
		subtotal
	FROM order_items WHERE order_id = ?`
	rows, err := database.DB.Query(query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orderItems []OrderItem
	for rows.Next() {
		var orderItem OrderItem
		err = rows.Scan(&orderItem.ID, &orderItem.OrderID, &orderItem.ProductID, &orderItem.VariantID,
			&orderItem.ProductName, &orderItem.Quantity, &orderItem.PriceSnapshot, &orderItem.Subtotal)
		if err != nil {
			return nil, err
		}

		orderItems = append(orderItems, orderItem)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &orderItems, nil
}

func GetTotalAmountFromOrderItems(orderID int64) (*string, error) {
	query := `SELECT SUM(subtotal) FROM order_items WHERE order_id = ?`
	row := database.DB.QueryRow(query, orderID)

	var subtotal string
	err := row.Scan(&subtotal)
	if err != nil {
		return nil, err
	}

	return &subtotal, nil
}
