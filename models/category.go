package models

import "github.com/BatJoz21/my-online-shop-go-api/database"

type Category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func GetCategories() (*[]Category, error) {
	query := `SELECT
		id,
		name,
		slug
	FROM categories`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}

	var categories []Category
	for rows.Next() {
		var category Category
		err = rows.Scan(&category.ID, &category.Name, &category.Slug)
		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	return &categories, nil
}
