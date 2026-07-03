package database

func seedCategoriesTable() {
	categories := map[string]string{
		"Apparel":     "apparel",
		"Footwear":    "footwear",
		"Accessories": "accessories",
	}

	for name, slug := range categories {
		_, err := DB.Exec(`INSERT IGNORE INTO categories(name, slug)
			VALUES (?, ?)
		`, name, slug)
		if err != nil {
			panic("Failed to populate categories table")
		}
	}
}
