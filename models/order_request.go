package models

type NewOrderDTO struct {
	ShippingAddress string `json:"shipping_address" binding:"required"`
}

type NewOrderItemDTO struct {
	OrderID       int64  `json:"order_id" binding:"required"`
	ProductID     string `json:"product_id" binding:"required"`
	VariantID     string `json:"variant_id" binding:"required"`
	ProductName   string `json:"product_name_snapshot" binding:"required"`
	Quantity      string `json:"quantity" binding:"required"`
	PriceSnapshot string `json:"price_snapshot" binding:"required"`
}
