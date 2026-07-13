package models

import (
	"time"

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
	CreatedAt       time.Time       `json:"created_at"`
}

func (o *Order) GenerateNew() error {
	query := `INSERT INTO orders(user_id, order_number, status, total_amount, shipping_address)
		VALUES (?, ?, ?, ?, ?)`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(o.UserID, o.OrderNumber, o.Status, 0, o.ShippingAddress)
	if err != nil {
		return err
	}

	o.ID, err = result.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

func GetOrders(userID int64) (*[]Order, error) {
	query := `SELECT
		id,
		user_id,
		order_number,
		status,
		total_amount,
		shipping_address,
		created_at
	FROM orders WHERE user_id = ?`
	rows, err := database.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		err := rows.Scan(&order.ID, &order.UserID, &order.OrderNumber, &order.Status,
			&order.TotalAmount, &order.ShippingAddress, &order.CreatedAt)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &orders, nil
}

func GetOrder(id int64) (*Order, error) {
	query := `SELECT
		id,
		user_id,
		order_number,
		status,
		total_amount,
		shipping_address
	FROM orders WHERE id = ?`
	row := database.DB.QueryRow(query, id)

	var order Order
	err := row.Scan(&order.ID, &order.UserID, &order.OrderNumber, &order.Status,
		&order.TotalAmount, &order.ShippingAddress)
	if err != nil {
		return nil, err
	}

	return &order, nil
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
	defer stmt.Close()

	_, err = stmt.Exec(o.OrderNumber, o.TotalAmount, o.ID)

	return err
}
