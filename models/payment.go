package models

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/BatJoz21/my-online-shop-go-api/database"
	"github.com/BatJoz21/my-online-shop-go-api/utils"
	"github.com/shopspring/decimal"
)

type Payment struct {
	ID            int64            `json:"id"`
	OrderID       int64            `json:"order_id"`
	Provider      string           `json:"provider"`
	TransactionID utils.NullString `json:"transaction_id"`
	Amount        decimal.Decimal  `json:"amount"`
	Status        string           `json:"status"`
	PaidAt        utils.NullTime   `json:"paid_at"`
	RawResponse   json.RawMessage  `json:"raw_response"`
	CreatedAt     time.Time        `json:"created_at"`
}

func (p *Payment) Save() error {
	query := `INSERT INTO payments(order_id, provider, amount, status, raw_response)
		VALUES (?, ?, ?, ?, ?)`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(p.OrderID, p.Provider, p.Amount, p.Status, p.RawResponse)
	if err != nil {
		return err
	}

	p.ID, err = result.LastInsertId()

	return err
}

func UpdateFromWebhook(orderNumber string, newStatus, transactionID string, rawResponse []byte) error {
	orderID, orderStatus, err := GetIDStatusOrder(orderNumber)
	if err != nil {
		return err
	}

	if orderStatus == "paid" && newStatus != "success" {
		return nil
	}

	var paidAt interface{}
	if newStatus == "success" {
		paidAt = time.Now()
	} else {
		paidAt = nil
	}

	query := `UPDATE payments SET
		status = ?,
		transaction_id = ?,
		paid_at = ?,
		raw_response = ?
	WHERE order_id = ?`

	_, err = database.DB.Exec(query, newStatus, transactionID, paidAt, rawResponse, orderID)
	if err != nil {
		return err
	}

	if newStatus == "success" && orderStatus == "pending" {
		if err := SetOrderToPaid(orderNumber); err != nil {
			return err
		}
	}

	return nil
}

func IsExistingPaymentPending(orderID int64) (bool, error) {
	query := `SELECT id FROM payments WHERE order_id = ? AND status = ?`
	row := database.DB.QueryRow(query, orderID, "pending")

	var id int64
	err := row.Scan(&id)
	if err == nil {
		return true, nil
	} else if err != sql.ErrNoRows {
		return true, err
	}

	return false, nil
}
