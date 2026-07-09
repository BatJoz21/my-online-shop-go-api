package models

type ProductVariantDTO struct {
	Name          string `json:"name" binding:"required"`
	Sku           string `json:"sku" binding:"required"`
	PriceModifier string `json:"price_modifier" binding:"required"`
	Stock         string `json:"stock" binding:"required"`
}

type UpdateStockDTO struct {
	Stock int64 `json:"stock" binding:"required"`
}
