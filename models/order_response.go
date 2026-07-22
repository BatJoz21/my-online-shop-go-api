package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type GetDetailedOrderDTO struct {
	ID               int64           `json:"id"`
	UserID           int64           `json:"user_id"`
	OrderNumber      string          `json:"order_number"`
	Status           string          `json:"status"`
	TotalAmount      decimal.Decimal `json:"total_amount"`
	ShippingAddress  string          `json:"shipping_address"`
	EstimatedArrival *time.Time      `json:"estimated_arrival"`
	CreatedAt        time.Time       `json:"created_at"`
	OwnerName        string          `json:"owner_name"`
}

type OrderForPayment struct {
	ID            int64           `json:"id"`
	OrderNumber   string          `json:"order_number"`
	Status        string          `json:"status"`
	TotalAmount   decimal.Decimal `json:"total_amount"`
	CustomerName  string          `json:"customer_name"`
	CustomerEmail string          `json:"customer_email"`
}
