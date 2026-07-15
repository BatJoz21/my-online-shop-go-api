package models

import "time"

type GetReview struct {
	ID        int64      `json:"id"`
	ProductID int64      `json:"product_id"`
	UserID    int64      `json:"user_id"`
	OrderID   int64      `json:"order_id"`
	Rating    int        `json:"rating"`
	Comment   string     `json:"comment"`
	CreatedAt *time.Time `json:"created_at"`
	UserName  string     `json:"username"`
}
