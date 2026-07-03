package database

func createTables() {
	createUsersTable := `CREATE TABLE IF NOT EXISTS users (
		id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(150) UNIQUE NOT NULL,
		password_hash VARCHAR(225) NOT NULL,
		role ENUM('customer', 'merchant', 'admin') DEFAULT 'customer',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	)`
	_, err := DB.Exec(createUsersTable)
	if err != nil {
		panic("Failed to create tasks table: " + err.Error())
	}

	createRefreshTokenTable := `CREATE TABLE IF NOT EXISTS refresh_tokens ( 
		id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT, 
		user_id BIGINT UNSIGNED NOT NULL, 
		device_name VARCHAR(100) NULL, 
		token_hash VARCHAR(225) NOT NULL, 
		expires_at DATETIME NOT NULL, 
		revoked_at DATETIME NULL, 
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, 
		
		CONSTRAINT token_user_id_fk 
			FOREIGN KEY(user_id) 
			REFERENCES users (id) 
			ON DELETE CASCADE 
			ON UPDATE CASCADE
	)`
	_, err = DB.Exec(createRefreshTokenTable)
	if err != nil {
		panic("Failed to create refresh_tokens table: " + err.Error())
	}

	createCategoriesTable := `CREATE TABLE IF NOT EXISTS categories (
		id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
		name VARCHAR(100) NOT NULL,
		slug VARCHAR(120) UNIQUE NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	)`
	_, err = DB.Exec(createCategoriesTable)
	if err != nil {
		panic("Failed to create categories table: " + err.Error())
	}

	createProductsTable := `CREATE TABLE IF NOT EXISTS products (
		id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
		category_id BIGINT UNSIGNED,
		name VARCHAR(180) NOT NULL,
		slug VARCHAR(180) UNIQUE NOT NULL,
		description TEXT,
		price DECIMAL(12, 2) NOT NULL,
		stock INT DEFAULT 0,
		image VARCHAR(225) NULL,
		is_active BOOLEAN DEFAULT TRUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

		CONSTRAINT products_category_id_fk 
			FOREIGN KEY(category_id) 
			REFERENCES categories (id) 
			ON DELETE CASCADE 
			ON UPDATE CASCADE
	)`
	_, err = DB.Exec(createProductsTable)
	if err != nil {
		panic("Failed to create products table: " + err.Error())
	}

	createCartsTable := `CREATE TABLE IF NOT EXISTS carts (
		id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
		user_id BIGINT UNSIGNED UNIQUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

		CONSTRAINT carts_user_id_fk 
			FOREIGN KEY(user_id) 
			REFERENCES users (id) 
			ON DELETE CASCADE 
			ON UPDATE CASCADE
	)`
	_, err = DB.Exec(createCartsTable)
	if err != nil {
		panic("Failed to create carts table: " + err.Error())
	}

	createCartItemsTable := `CREATE TABLE IF NOT EXISTS cart_items (
		id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
		cart_id BIGINT UNSIGNED,
		product_id BIGINT UNSIGNED,
		quantity INT,
		price_snapshot DECIMAL(12, 2),
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

		CONSTRAINT cartitem_cart_id_fk 
			FOREIGN KEY(cart_id) 
			REFERENCES carts (id) 
			ON DELETE CASCADE 
			ON UPDATE CASCADE,

		CONSTRAINT cartitem_product_id_fk 
			FOREIGN KEY(product_id) 
			REFERENCES products (id) 
			ON DELETE CASCADE 
			ON UPDATE CASCADE
	)`
	_, err = DB.Exec(createCartItemsTable)
	if err != nil {
		panic("Failed to create cart_items table: " + err.Error())
	}

	createOrdersTable := `CREATE TABLE IF NOT EXISTS orders (
		id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
		user_id BIGINT UNSIGNED,
		order_number VARCHAR(30) UNIQUE,
		status ENUM('pending', 'paid', 'shipped', 'completed', 'cancelled') DEFAULT 'pending',
		total_amount DECIMAL(12, 2),
		shipping_address TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

		CONSTRAINT order_user_id_fk 
			FOREIGN KEY(user_id) 
			REFERENCES users (id) 
			ON DELETE CASCADE 
			ON UPDATE CASCADE
	)`
	_, err = DB.Exec(createOrdersTable)
	if err != nil {
		panic("Failed to create orders table: " + err.Error())
	}

	createOrderItemsTable := `CREATE TABLE IF NOT EXISTS order_items (
		id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
		order_id BIGINT UNSIGNED,
		product_id BIGINT UNSIGNED,
		product_name_snapshot VARCHAR(150),
		quantity INT,
		price_shapshot DECIMAL(12, 2),
		subtotal DECIMAL(12, 2),
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

		CONSTRAINT orderitem_order_id_fk 
			FOREIGN KEY(order_id) 
			REFERENCES orders (id) 
			ON DELETE CASCADE 
			ON UPDATE CASCADE,

		CONSTRAINT orderitem_product_id_fk 
			FOREIGN KEY(product_id) 
			REFERENCES products (id) 
			ON DELETE CASCADE 
			ON UPDATE CASCADE
	)`
	_, err = DB.Exec(createOrderItemsTable)
	if err != nil {
		panic("Failed to create order_items table: " + err.Error())
	}

	createPaymentsTable := `CREATE TABLE IF NOT EXISTS payments (
		id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
		order_id BIGINT UNSIGNED,
		provider VARCHAR(50),
		transaction_id VARCHAR(150),
		amount DECIMAL(12, 2),
		status ENUM('pending', 'success', 'failed', 'expired') DEFAULT 'pending',
		paid_at DATETIME NULL,
		raw_response JSON NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

		CONSTRAINT payment_order_id_fk 
			FOREIGN KEY(order_id) 
			REFERENCES orders (id) 
			ON DELETE CASCADE 
			ON UPDATE CASCADE
	)`
	_, err = DB.Exec(createPaymentsTable)
	if err != nil {
		panic("Failed to create payments table: " + err.Error())
	}

	createReviewsTable := `CREATE TABLE IF NOT EXISTS reviews (
		id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
		product_id BIGINT UNSIGNED,
		user_id BIGINT UNSIGNED,
		order_id BIGINT UNSIGNED,
		rating TINYINT NOT NULL,
		comment TEXT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

		CONSTRAINT review_product_id_fk 
			FOREIGN KEY(product_id) 
			REFERENCES products (id) 
			ON DELETE CASCADE 
			ON UPDATE CASCADE,

		CONSTRAINT review_user_id_fk 
			FOREIGN KEY(user_id) 
			REFERENCES users (id) 
			ON DELETE CASCADE 
			ON UPDATE CASCADE,

		CONSTRAINT review_order_id_fk 
			FOREIGN KEY(order_id) 
			REFERENCES orders (id) 
			ON DELETE CASCADE 
			ON UPDATE CASCADE,

		CONSTRAINT chk_rating_range CHECK (rating BETWEEN 1 AND 5)
	)`
	_, err = DB.Exec(createReviewsTable)
	if err != nil {
		panic("Failed to create cart_items table: " + err.Error())
	}
}
