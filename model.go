package main

import (
	"database/sql"
)

type product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func (prod *product) getProduct(db *sql.DB) error {
	return db.QueryRow("SELECT FROM product WHERE id = $1", prod.ID).Scan(&prod.Name, &prod.Price)
}

func (prod *product) updateProduct(db *sql.DB) error {
	_, err := db.Exec("UPDATE product SET name = $1, price = $2 WHERE id = $3",
		prod.Name,
		prod.Price,
		prod.ID)

	return err
}

func (prod *product) createProduct(db *sql.DB) error {

	return db.QueryRow("INSERT INTO product(name, price) VALUES($1, $2) RETURNING id",
		prod.Name,
		prod.Price).Scan(&prod.ID)
}

func getAllProducts(db *sql.DB, start, count int) ([]product, error) {

	rows, err := db.Query("SELECT * FROM product LIMIT $1 OFFSET $2", count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	products := []product{}

	for rows.Next() {

		var p product

		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			return nil, err
		}

		products = append(products, p)
	}

	return products, nil
}
