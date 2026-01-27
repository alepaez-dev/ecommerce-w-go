package orders

import (
	"context"
	"errors"
	"fmt"

	repo "github.com/alepaez-dev/ecommerce/internal/adapters/postgresql/sqlc"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrProductNoStock  = errors.New("product has not enough stock")
	ErrRequiredValue   = errors.New("is required")
)

type Service interface {
	PlaceOrder(ctx context.Context, order createOrderParams) (repo.Order, error)
}

type svc struct {
	txManager TxManager
	products  ProductStoreFactory
}

func NewService(txManager TxManager, products ProductStoreFactory) Service {
	return &svc{txManager: txManager, products: products}
}

func (s *svc) PlaceOrder(ctx context.Context, tempOrder createOrderParams) (repo.Order, error) {
	if tempOrder.CustomerID == 0 {
		return repo.Order{}, fmt.Errorf("%w customer ID is required", ErrRequiredValue)
	}
	if len(tempOrder.Items) == 0 {
		return repo.Order{}, fmt.Errorf("at least one item %w", ErrRequiredValue)
	}

	var createdOrder repo.Order
	err := s.txManager.WithTx(ctx, func(q repo.Querier) error {
		productStore := s.products(q)

		order, err := q.CreateOrder(ctx, tempOrder.CustomerID)
		if err != nil {
			return err
		}

		for _, item := range tempOrder.Items {
			product, err := productStore.FindProduct(ctx, item.ProductID)
			if err != nil {
				return fmt.Errorf("%w (product_id=%d)", ErrProductNotFound, item.ProductID)
			}

			if product.Quantity < item.Quantity {
				return fmt.Errorf("%w (product_id=%d)", ErrProductNoStock, item.ProductID)
			}

			if _, err := productStore.DecrementStock(ctx, item.ProductID, item.Quantity); err != nil {
				return err
			}

			if _, err := q.CreateOrderItem(ctx, repo.CreateOrderItemParams{
				OrderID:      order.ID,
				ProductID:    item.ProductID,
				Quantity:     item.Quantity,
				PriceInCents: product.PriceInCents,
			}); err != nil {
				return err
			}
		}

		createdOrder = order
		return nil
	})

	if err != nil {
		return repo.Order{}, err
	}

	return createdOrder, nil
}
