package repository

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/DarkReduX/productservice/model"
)

type Product struct {
	pool *pgxpool.Pool
}

func NewProduct(pool *pgxpool.Pool) *Product {
	return &Product{pool: pool}
}

func (r *Product) Create(ctx context.Context, product *model.Product) (*model.Product, error) {
	product.ID = uuid.New()

	_, err := r.pool.Exec(ctx, `INSERT INTO products (id, user_id,name, price, description) VALUES ($1, $2, $3, $4, $5)`,
		product.ID, product.UserID, product.Name, product.Price, product.Description)

	return product, err
}

func (r *Product) Get(ctx context.Context, id uuid.UUID) (*model.Product, error) {
	product := new(model.Product)

	err := pgxscan.Get(ctx, r.pool, product, `SELECT id, user_id, name, price, description FROM products WHERE id=$1`, id)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *Product) List(ctx context.Context) ([]*model.Product, error) {
	var products []*model.Product

	err := pgxscan.Select(ctx, r.pool, &products, `SELECT id, user_id, name, price, description FROM products`)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (r *Product) ListByUser(ctx context.Context, userID uuid.UUID) ([]*model.Product, error) {
	var products []*model.Product

	err := pgxscan.Select(
		ctx, r.pool, &products,
		`SELECT id, user_id, name, price, description FROM products WHERE user_id=$1`,
		userID,
	)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (r *Product) Update(ctx context.Context, product *model.Product) error {
	_, err := r.pool.Exec(ctx, `UPDATE products SET name=$1, price=$2, description=$3 WHERE id=$4 AND user_id=$5`,
		product.Name, product.Price, product.Description, product.ID, product.UserID)

	return err
}

func (r *Product) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM products WHERE id=$1`, id)

	return err
}
