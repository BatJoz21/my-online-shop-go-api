package models

import "github.com/BatJoz21/my-online-shop-go-api/database"

type Cart struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"user_id"`
}

func GenerateNewCart(userId int64) error {
	query := `INSERT INTO carts(user_id) VALUES(?)`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userId)
	if err != nil {
		return err
	}

	return nil
}

func GetUserCartID(userID int64) (*int64, error) {
	query := `SELECT id FROM carts WHERE user_id = ?`
	row := database.DB.QueryRow(query, userID)

	var id int64
	err := row.Scan(&id)
	if err != nil {
		return nil, err
	}

	return &id, nil
}
