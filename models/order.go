package models

import (
	"time"

	"github.com/BatJoz21/my-online-shop-go-api/database"
	"github.com/shopspring/decimal"
)

type Order struct {
	ID               int64           `json:"id"`
	UserID           int64           `json:"user_id"`
	OrderNumber      string          `json:"order_number"`
	Status           string          `json:"status"`
	TotalAmount      decimal.Decimal `json:"total_amount"`
	ShippingAddress  string          `json:"shipping_address"`
	EstimatedArrival *time.Time      `json:"estimated_arrival"`
	CreatedAt        time.Time       `json:"created_at"`
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

func GetAllOrders(status string) (*[]GetDetailedOrderDTO, error) {
	query := `SELECT
		orders.id,
		orders.user_id,
		orders.order_number,
		orders.status,
		orders.total_amount,
		orders.shipping_address,
		orders.estimated_arrival,
		orders.created_at,
		users.name as owner_name
	FROM orders
	JOIN users ON orders.user_id = users.id`

	var args []any
	if status != "" {
		query += ` WHERE orders.status = ?`
		args = append(args, status)
	}

	query += ` ORDER BY orders.id ASC`

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []GetDetailedOrderDTO
	for rows.Next() {
		var order GetDetailedOrderDTO
		err := rows.Scan(&order.ID, &order.UserID, &order.OrderNumber, &order.Status,
			&order.TotalAmount, &order.ShippingAddress, &order.EstimatedArrival,
			&order.CreatedAt, &order.OwnerName)
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

func GetOrders(userID int64, status string) (*[]Order, error) {
	query := `SELECT
		id,
		user_id,
		order_number,
		status,
		total_amount,
		shipping_address,
		estimated_arrival,
		created_at
	FROM orders 
	WHERE user_id = ?`

	var args []any
	args = append(args, userID)
	if status != "" {
		query += ` AND status = ?`
		args = append(args, status)
	}

	query += ` ORDER BY id ASC`

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		err := rows.Scan(&order.ID, &order.UserID, &order.OrderNumber, &order.Status,
			&order.TotalAmount, &order.ShippingAddress, &order.EstimatedArrival, &order.CreatedAt)
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
		shipping_address,
		estimated_arrival,
		created_at
	FROM orders WHERE id = ?`
	row := database.DB.QueryRow(query, id)

	var order Order
	err := row.Scan(&order.ID, &order.UserID, &order.OrderNumber, &order.Status,
		&order.TotalAmount, &order.ShippingAddress, &order.EstimatedArrival, &order.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func GetOrderForMerchant(id int64) (*GetDetailedOrderDTO, error) {
	query := `SELECT
		orders.id,
		orders.user_id,
		orders.order_number,
		orders.status,
		orders.total_amount,
		orders.shipping_address,
		orders.estimated_arrival,
		orders.created_at,
		users.name as owner_name
	FROM orders
	JOIN users ON orders.user_id = users.id
	WHERE orders.id = ?`
	row := database.DB.QueryRow(query, id)

	var order GetDetailedOrderDTO
	err := row.Scan(&order.ID, &order.UserID, &order.OrderNumber, &order.Status,
		&order.TotalAmount, &order.ShippingAddress, &order.EstimatedArrival,
		&order.CreatedAt, &order.OwnerName)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func GetOrderForShowPage(id int64) (*Order, *[]OrderItem, error) {
	order, err := GetOrder(id)
	if err != nil {
		return nil, nil, err
	}

	orderItems, err := GetAllItemFromOrder(id)
	if err != nil {
		return nil, nil, err
	}

	return order, orderItems, nil
}

func GetOrderForMerchantShowPage(id int64) (*GetDetailedOrderDTO, *[]OrderItem, error) {
	order, err := GetOrderForMerchant(id)
	if err != nil {
		return nil, nil, err
	}

	orderItems, err := GetAllItemFromOrder(id)
	if err != nil {
		return nil, nil, err
	}

	return order, orderItems, nil
}

func GetOrderForPayment(orderID int64) (*OrderForPayment, error) {
	query := `SELECT
		orders.id,
		orders.order_number,
		orders.status,
		orders.total_amount,
		users.name,
		users.email
	FROM orders
	JOIN users ON orders.user_id = users.id
	WHERE orders.id = ?`
	row := database.DB.QueryRow(query, orderID)

	var order OrderForPayment
	err := row.Scan(&order.ID, &order.OrderNumber, &order.Status,
		&order.TotalAmount, &order.CustomerName, &order.CustomerEmail)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func GetIDStatusOrder(orderNumber string) (int64, string, error) {
	query := `SELECT
		id,
		status
	FROM orders
	WHERE order_number = ?`
	row := database.DB.QueryRow(query, orderNumber)

	var id int64
	var status string
	err := row.Scan(&id, &status)
	if err != nil {
		return -1, "", err
	}

	return id, status, nil
}

func IsOrderComplete(id int64) (bool, error) {
	query := `SELECT status FROM orders WHERE id = ?`
	row := database.DB.QueryRow(query, id)

	var status string
	err := row.Scan(&status)
	if err != nil {
		return false, nil
	}

	if status == "completed" {
		return true, nil
	}

	return false, nil
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

func (o *Order) Update() error {
	query := `UPDATE orders SET
		order_number = ?,
		status = ?,
		total_amount = ?,
		shipping_address = ?,
		estimated_arrival = ?
	WHERE id = ?`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(o.OrderNumber, o.Status, o.TotalAmount,
		o.ShippingAddress, o.EstimatedArrival, o.ID)

	return err
}

func SetOrderToPaid(orderNumber string) error {
	query := `UPDATE orders SET
		status = ?
	WHERE order_number = ?`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec("paid", orderNumber)

	return err
}

func (o *Order) Delete() error {
	query := `DELETE FROM orders WHERE id = ?`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(o.ID)

	return err
}
