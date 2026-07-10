package models

import (
	"github.com/BatJoz21/my-online-shop-go-api/database"
	"github.com/shopspring/decimal"
)

const (
	ProductPerPageLimit = 5
)

type Product struct {
	ID           int64           `json:"id"`
	CategoryID   int64           `json:"category_id"`
	Name         string          `json:"name"`
	Slug         string          `json:"slug"`
	Description  string          `json:"description"`
	Price        decimal.Decimal `json:"price"`
	Image        *string         `json:"image"`
	IsActive     bool            `json:"is_active"`
	CategoryName string          `json:"category_name"`
}

func (p *Product) Save() error {
	query := `INSERT INTO products(category_id, name, slug, description, price, image)
		VALUES (?, ?, ?, ?, ?, ?)`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(p.CategoryID, p.Name, p.Slug, p.Description, p.Price, p.Image)
	if err != nil {
		return err
	}

	p.ID, err = result.LastInsertId()
	p.IsActive = true
	if err != nil {
		return err
	}

	return nil
}

func GetAllProducts(category string, offset int) ([]Product, error) {
	query := `SELECT
		products.id,
		products.category_id,
		products.name,
		products.slug,
		products.description,
		products.price,
		products.image,
		products.is_active,
		categories.name AS category_name
	FROM products
	JOIN categories ON products.category_id = categories.id
	WHERE products.is_active = ?
	ORDER BY products.id ASC
	LIMIT ? OFFSET ?`

	rows, err := database.DB.Query(query, 1, ProductPerPageLimit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		err = rows.Scan(&product.ID, &product.CategoryID, &product.Name, &product.Slug,
			&product.Description, &product.Price, &product.Image,
			&product.IsActive, &product.CategoryName)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func GetAllStockedProducts(category int64, search string, offset int) (*[]Product, error) {
	query := `SELECT
		products.id,
		products.category_id,
		products.name,
		products.slug,
		products.description,
		products.price,
		products.image,
		products.is_active,
		categories.name AS category_name
	FROM products
	JOIN product_variants ON products.id = product_variants.product_id
	JOIN categories ON products.category_id = categories.id
	WHERE products.is_active = ?`

	var args []any
	args = append(args, 1)

	if category != 0 {
		query += ` AND products.category_id = ?`
		args = append(args, category)
	}

	if search != "" {
		query += ` AND products.name LIKE ?`
		search = "%" + search + "%"
		args = append(args, search)
	}

	query += ` 
	GROUP BY product_variants.product_id
	HAVING SUM(product_variants.stock) > 0
	ORDER BY products.id ASC
	LIMIT ? OFFSET ?`

	args = append(args, ProductPerPageLimit)
	args = append(args, offset)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		err = rows.Scan(&product.ID, &product.CategoryID, &product.Name, &product.Slug,
			&product.Description, &product.Price, &product.Image,
			&product.IsActive, &product.CategoryName)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &products, nil
}

func GetActiveProduct(id int64) (*Product, error) {
	query := `SELECT
		products.id,
		products.category_id,
		products.name,
		products.slug,
		products.description,
		products.price,
		products.image,
		products.is_active,
		categories.name AS category_name
	FROM products
	JOIN product_variants ON products.id = product_variants.product_id
	JOIN categories ON products.category_id = categories.id
	WHERE products.id = ? AND is_active = ?
	GROUP BY product_variants.product_id
	HAVING SUM(product_variants.stock) > 0`
	row := database.DB.QueryRow(query, id, 1)

	var product Product
	err := row.Scan(&product.ID, &product.CategoryID, &product.Name, &product.Slug,
		&product.Description, &product.Price, &product.Image,
		&product.IsActive, &product.CategoryName)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func GetProduct(id int64) (*Product, error) {
	query := `SELECT
		products.id,
		products.category_id,
		products.name,
		products.slug,
		products.description,
		products.price,
		products.image,
		products.is_active,
		categories.name AS category_name
	FROM products
	JOIN categories ON products.category_id = categories.id
	WHERE products.id = ?`
	row := database.DB.QueryRow(query, id)

	var product Product
	err := row.Scan(&product.ID, &product.CategoryID, &product.Name, &product.Slug,
		&product.Description, &product.Price, &product.Image,
		&product.IsActive, &product.CategoryName)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (p *Product) Update() error {
	query := `UPDATE products SET
		category_id = ?,
		name = ?,
		slug = ?,
		description = ?,
		price = ?,
		image = ?
	WHERE id = ? AND is_active = ?`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(p.CategoryID, p.Name, p.Slug, p.Description, p.Price, p.Image, p.ID, 1)

	return err
}

func (p *Product) Restore() error {
	query := `UPDATE products SET is_active = ? WHERE id = ?`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(1, p.ID)

	return err
}

func (p *Product) SoftDelete() error {
	query := `UPDATE products SET is_active = 0 WHERE id = ?`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(p.ID)

	return err
}

func (p *Product) Delete() error {
	query := `DELETE FROM products WHERE id = ?`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(p.ID)

	return err
}
