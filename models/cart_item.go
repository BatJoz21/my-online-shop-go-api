package models

import (
	"github.com/BatJoz21/my-online-shop-go-api/database"
	"github.com/shopspring/decimal"
)

type CartItems struct {
	ID            int64           `json:"id"`
	CartID        int64           `json:"cart_id"`
	ProductID     int64           `json:"product_id"`
	VariantID     int64           `json:"variant_id"`
	Quantity      int             `json:"quantity"`
	PriceSnapshot decimal.Decimal `json:"price_snapshot"`
}

func (c *CartItems) Save() error {
	query := `INSERT INTO cart_items(cart_id, product_id, variant_id, quantity, price_snapshot)
		VALUES (?, ?, ?, ?, ?)`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(c.CartID, c.ProductID, c.VariantID, c.Quantity, c.PriceSnapshot)
	if err != nil {
		return err
	}

	c.ID, err = result.LastInsertId()
	return err
}

func GetAllItemInCart(cartID int64) (*[]CartItemResponse, error) {
	query := `SELECT
		cart_items.id,
		cart_items.cart_id,
		cart_items.product_id,
		products.name as product_name,
		cart_items.variant_id,
		product_variants.name as variant_name,
		cart_items.quantity,
		cart_items.price_snapshot
	FROM cart_items
	JOIN products ON cart_items.product_id = products.id
	JOIN product_variants ON cart_items.variant_id = product_variants.id
	WHERE cart_items.cart_id = ?
	ORDER BY products.name ASC`
	rows, err := database.DB.Query(query, cartID)
	if err != nil {
		return nil, err
	}

	var cartItems []CartItemResponse
	for rows.Next() {
		var cartItem CartItemResponse
		err = rows.Scan(&cartItem.ID, &cartItem.CartID, &cartItem.ProductID,
			&cartItem.ProductName, &cartItem.VariantID, &cartItem.VariantName,
			&cartItem.Quantity, &cartItem.PriceSnapshot)
		if err != nil {
			return nil, err
		}

		cartItems = append(cartItems, cartItem)
	}

	if err := rows.Err(); err != nil {
		return nil, err 
	}

	return &cartItems, nil
}

func GetTotalItemInACart(cartID int64) (int, error) {
	query := `SELECT COUNT(*) FROM cart_items WHERE cart_id = ?`
	row := database.DB.QueryRow(query, cartID)

	var total int
	err := row.Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (c *CartItems) Update() error {
	query := `UPDATE cart_items SET
		variant_id = ?,
		quantity = ?,
		price_snapshot = ?
	WHERE id = ?`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(c.VariantID, c.Quantity, c.PriceSnapshot, c.ID)
	if err != nil {
		return err
	}

	return nil
}

func DeleteCartItem(id int64) error {
	query := `DELETE FROM cart_items WHERE id = ?`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)

	return err
}

func EmptyCart(cartID int64) error {
	query := `DELETE FROM cart_items WHERE cart_id = ?`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(cartID)

	return err
}
