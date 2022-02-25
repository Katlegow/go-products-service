package main

import (
	"database/sql"
	"errors"
)

type product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func (prod *product) getProduct(db *sql.DB) error {
	return errors.New("Not implemented yet!")
}

func (prod *product) updateProduct(db *sql.DB) error {
	return errors.New("Not implemented yet!")
}

func (prod *product) createProduct(db *sql.DB) error {
	return errors.New("Not implemented yet!")
}

func getAllProducts(db *sql.DB, start, count int) ([]product, error) {
	return nil, errors.New("Not implemented yet!")
}
