package models

type NewOrderDTO struct {
	TotalAmount     string `json:"total_amount" binding:"required"`
	ShippingAddress string `json:"shipping_address" binding:"required"`
}
