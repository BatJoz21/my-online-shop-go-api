package models

import "github.com/BatJoz21/my-online-shop-go-api/database"

func GetDashboardStats() (*DashboardStats, error) {
	query := `SELECT
		(SELECT COUNT(*) FROM products WHERE is_active = 1) as total_products,
		(SELECT COUNT(*) FROM orders WHERE status = "pending") as pending_orders,
		(SELECT SUM(total_amount) FROM orders WHERE status = "completed") as total_revenue,
		(SELECT COUNT(*) FROM product_variants WHERE stock < 5) as low_stock_count`
	row := database.DB.QueryRow(query)

	var stats DashboardStats
	err := row.Scan(&stats.TotalProducts, &stats.PendingOrders, &stats.TotalRevenue, &stats.LowStockCount)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

func GetRecentOrders() (*[]RecentOrderDTO, error) {
	query := `SELECT id, order_number, status, total_amount FROM orders ORDER BY created_at DESC LIMIT 5`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []RecentOrderDTO
	for rows.Next() {
		var order RecentOrderDTO
		err = rows.Scan(&order.ID, &order.OrderNumber, &order.Status, &order.TotalAmount)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &orders, nil
}

func GetLowStockProducts() (*[]LowStockProductDTO, error) {
	query := `SELECT
		products.name as product_name,
		product_variants.name as variant_name,
		product_variants.stock
	FROM product_variants
	JOIN products ON product_variants.product_id = products.id
	WHERE product_variants.stock < 5`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []LowStockProductDTO
	for rows.Next() {
		var dto LowStockProductDTO
		err = rows.Scan(&dto.ProductName, &dto.VariantName, &dto.Stock)
		if err != nil {
			return nil, err
		}

		products = append(products, dto)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &products, nil
}

func GetRecentReviews() (*[]RecentReviewDTO, error) {
	query := `SELECT
		products.name as product_name,
		reviews.rating,
		reviews.comment
	FROM reviews
	JOIN products ON reviews.product_id = products.id
	ORDER BY reviews.created_at DESC`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []RecentReviewDTO
	for rows.Next() {
		var dto RecentReviewDTO
		err = rows.Scan(&dto.ProductName, &dto.Rating, &dto.Comment)
		if err != nil {
			return nil, err
		}

		reviews = append(reviews, dto)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &reviews, nil
}
