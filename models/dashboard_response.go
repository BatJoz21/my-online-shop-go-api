package models

import "github.com/shopspring/decimal"

type DashboardStats struct {
	TotalProducts int             `json:"total_products"`
	PendingOrders int             `json:"pending_orders"`
	TotalRevenue  decimal.Decimal `json:"total_revenue"`
	LowStockCount int             `json:"low_stock_count"`
}

type RecentOrderDTO struct {
	ID          int64           `json:"id"`
	OrderNumber string          `json:"order_number"`
	Status      string          `json:"status"`
	TotalAmount decimal.Decimal `json:"total_amount"`
}

type LowStockProductDTO struct {
	ProductName string `json:"product_name"`
	VariantName string `json:"variant_name"`
	Stock       int    `json:"stock"`
}

type RecentReviewDTO struct {
	ProductName string `json:"product_name"`
	Rating      int    `json:"rating"`
	Comment     string `json:"comment"`
}
