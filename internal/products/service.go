package products

import (
	"context"

	repo "github.com/alepaez-dev/ecommerce/internal/adapters/postgresql/sqlc"
)

type Service interface {
	ListProducts(ctx context.Context) ([]repo.Product, error)
	FindProduct(ctx context.Context, id int64) (repo.Product, error)
	DecrementStock(ctx context.Context, id int64, quantity int32) (repo.DecrementProductStockRow, error)
}

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &svc{repo: repo}
}

func (s *svc) ListProducts(ctx context.Context) ([]repo.Product, error) {
	return s.repo.ListProducts(ctx)
}

func (s *svc) FindProduct(ctx context.Context, id int64) (repo.Product, error) {
	return s.repo.FindProductByID(ctx, id)
}

func (s *svc) DecrementStock(ctx context.Context, id int64, quantity int32) (repo.DecrementProductStockRow, error) {
	return s.repo.DecrementProductStock(ctx, repo.DecrementProductStockParams{
		ID:       id,
		Quantity: quantity,
	})

}
