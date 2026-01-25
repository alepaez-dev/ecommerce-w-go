package orders

import (
	"context"

	repo "github.com/alepaez-dev/ecommerce/internal/adapters/postgresql/sqlc"
)

type orderItem struct {
	ProductID int64 `json:"product_id"`
	Quantity  int32 `json:"quantity"`
}

type createOrderParams struct {
	CustomerID int64       `json:"customer_id"`
	Items      []orderItem `json:"items"`
}

type ProductStore interface {
	FindProduct(ctx context.Context, id int64) (repo.Product, error)
	DecrementStock(ctx context.Context, id int64, quantity int32) (repo.DecrementProductStockRow, error)
}

type ProductStoreFactory func(repo.Querier) ProductStore

type TxManager interface {
	WithTx(ctx context.Context, fn func(repo.Querier) error) error
}
