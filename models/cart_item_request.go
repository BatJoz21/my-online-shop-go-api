package models

type CartItemsDTO struct {
	ProductID     int64  `json:"product_id" binding:"required"`
	VariantID     int64  `json:"variant_id" binding:"required"`
	Quantity      int    `json:"quantity" binding:"required"`
	PriceSnapshot string `json:"price_snapshot" binding:"required"`
}

type UpdateCartItemDTO struct {
	VariantID     int64  `json:"variant_id" binding:"required"`
	Quantity      int    `json:"quantity" binding:"required"`
	PriceSnapshot string `json:"price_snapshot" binding:"required"`
}
