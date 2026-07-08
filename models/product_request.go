package models

import "github.com/shopspring/decimal"

type ProductVariantDTO struct {
	Name          string          `json:"name" binding:"required"`
	Sku           string          `json:"sku" binding:"required"`
	PriceModifier decimal.Decimal `json:"price_modifier" binding:"required"`
	Stock         int64           `json:"stock" binding:"required"`
}

type UpdateStockDTO struct {
	Stock int64 `json:"stock" binding:"required"`
}
