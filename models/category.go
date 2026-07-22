package models

import "github.com/BatJoz21/my-online-shop-go-api/database"

type Category struct {
	ID   int64  `json:"id"`
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug" binding:"required"`
}

func (c *Category) Save() error {
	query := `INSERT INTO categories(name, slug) VALUES (?, ?)`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(c.Name, c.Slug)
	if err != nil {
		return err
	}

	c.ID, err = result.LastInsertId()
	if err != nil {
		return err
	}

	return nil
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

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &categories, nil
}
