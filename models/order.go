package models

import (
	"github.com/BatJoz21/my-online-shop-go-api/database"
	"github.com/shopspring/decimal"
)

type Order struct {
	ID              int64           `json:"id"`
	UserID          int64           `json:"user_id"`
	OrderNumber     string          `json:"order_number"`
	Status          string          `json:"status"`
	TotalAmount     decimal.Decimal `json:"total_amount"`
	ShippingAddress string          `json:"shipping_address"`
}

func (o *Order) GenerateNew() error {
	query := `INSERT INTO orders(user_id, order_number, status, total_amount, shipping_address)
		VALUES (?, ?, ?, ?, ?)`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(o.UserID, o.OrderNumber, o.Status,
		o.TotalAmount, o.ShippingAddress)
	if err != nil {
		return err
	}

	o.ID, err = result.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

func GetTodayTotalOrder() (int, error) {
	query := `SELECT COUNT(*) FROM orders WHERE created_at = CURDATE()`
	row := database.DB.QueryRow(query)

	var total int
	err := row.Scan(&total)
	if err != nil {
		return -1, err
	}

	return total, nil
}

func (o *Order) InputData() error {
	query := `UPDATE orders SET
		order_number = ?,
		total_amount = ?
	WHERE id = ?`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(o.OrderNumber, o.TotalAmount)

	return err
}
