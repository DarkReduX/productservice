package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/DarkReduX/productservice/internal/repository"
	"github.com/DarkReduX/productservice/model"
)

type Product struct {
	productRepo *repository.Product
}

func NewProduct(productRepo *repository.Product) *Product {
	return &Product{productRepo: productRepo}
}

func (s *Product) Create(ctx context.Context, product *model.Product) (*model.Product, error) {
	product.ID = uuid.New()
	return s.productRepo.Create(ctx, product)
}

func (s *Product) Get(ctx context.Context, id uuid.UUID) (*model.Product, error) {
	return s.productRepo.Get(ctx, id)
}

func (s *Product) List(ctx context.Context) ([]*model.Product, error) {
	return s.productRepo.List(ctx)
}

func (s *Product) ListByUser(ctx context.Context, userID uuid.UUID) ([]*model.Product, error) {
	return s.productRepo.ListByUser(ctx, userID)
}

func (s *Product) Update(ctx context.Context, product *model.Product) error {
	return s.productRepo.Update(ctx, product)
}

func (s *Product) Delete(ctx context.Context, id uuid.UUID) error {
	return s.productRepo.Delete(ctx, id)
}
