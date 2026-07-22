package models

import (
	"time"

	"github.com/BatJoz21/my-online-shop-go-api/database"
)

type Review struct {
	ID        int64      `json:"id"`
	ProductID int64      `json:"product_id"`
	UserID    int64      `json:"user_id"`
	OrderID   int64      `json:"order_id"`
	Rating    int        `json:"rating"`
	Comment   string     `json:"comment"`
	CreatedAt *time.Time `json:"created_at"`
}

func (r *Review) Save() error {
	query := `INSERT INTO reviews(
		product_id, 
		user_id, 
		order_id, 
		rating,
		comment)
	VALUES (?, ?, ?, ?, ?)`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(r.ProductID, r.UserID, r.OrderID,
		r.Rating, r.Comment)
	if err != nil {
		return err
	}

	r.ID, err = result.LastInsertId()
	return err
}

func GetProductReviews(productID int64) (*[]GetReview, error) {
	query := `SELECT
		reviews.id,
		reviews.product_id, 
		reviews.user_id, 
		reviews.order_id, 
		reviews.rating,
		reviews.comment,
		reviews.created_at,
		users.name
	FROM reviews
	JOIN users ON reviews.user_id = users.id
	WHERE reviews.product_id = ?`
	rows, err := database.DB.Query(query, productID)
	if err != nil {
		return nil, err
	}

	var reviews []GetReview
	for rows.Next() {
		var review GetReview
		err = rows.Scan(&review.ID, &review.ProductID, &review.UserID,
			&review.OrderID, &review.Rating, &review.Comment, &review.CreatedAt, &review.UserName)
		if err != nil {
			return nil, err
		}

		reviews = append(reviews, review)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &reviews, nil
}
