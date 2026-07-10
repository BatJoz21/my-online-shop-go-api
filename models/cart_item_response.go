package models

import "github.com/shopspring/decimal"

type CartItemResponse struct {
	ID            int64           `json:"id"`
	CartID        int64           `json:"cart_id"`
	ProductID     int64           `json:"product_id"`
	ProductName   string          `json:"product_name"`
	VariantID     int64           `json:"variant_id"`
	VariantName   string          `json:"variant_name"`
	Quantity      int             `json:"quantity"`
	PriceSnapshot decimal.Decimal `json:"price_snapshot"`
}
