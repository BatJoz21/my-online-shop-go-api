package models

import "github.com/shopspring/decimal"

type NewProductVariant struct {
	Name          string          `json:"name" binding:"required"`
	Sku           string          `json:"sku" binding:"required"`
	PriceModifier decimal.Decimal `json:"price_modifier" binding:"required"`
}
