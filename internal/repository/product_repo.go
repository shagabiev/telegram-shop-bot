package repository

import (
	"database/sql"

	"github.com/shagabiev/telegram-shop-bot/internal/models"
)

type ProductRepo struct {
	db *sql.DB
}

func NewProductRepo(db *sql.DB) *ProductRepo {
	return &ProductRepo{db: db}
}

func (r *ProductRepo) GetAll() ([]models.Product, error) {
	rows, err := r.db.Query(`SELECT id, name, description, price, quantity, photo_url FROM products`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Quantity, &p.PhotoURL)
		products = append(products, p)
	}
	return products, nil
}

func (r *ProductRepo) Add(p models.Product) error {
	_, err := r.db.Exec(`INSERT INTO products (name, description, price, quantity, photo_url) VALUES ($1,$2,$3,$4,$5)`,
		p.Name, p.Description, p.Price, p.Quantity, p.PhotoURL)
	return err
}

func (r *ProductRepo) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM products WHERE id = $1`, id)
	return err
}
