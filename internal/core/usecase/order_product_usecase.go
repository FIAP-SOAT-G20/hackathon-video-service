package usecase

import (
	"context"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/dto"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/port"
)

type orderProductUseCase struct {
	gateway port.OrderProductGateway
}

// NewOrderProductUseCase creates a new ListOrderProductsUseCase
func NewOrderProductUseCase(gateway port.OrderProductGateway) port.OrderProductUseCase {
	return &orderProductUseCase{gateway}
}

// List lists all orderProducts
func (uc *orderProductUseCase) List(ctx context.Context, i dto.ListOrderProductsInput) ([]*entity.OrderProduct, int64, error) {
	orderProducts, total, err := uc.gateway.FindAll(ctx, i.OrderID, i.ProductID, i.Page, i.Limit)
	if err != nil {
		return nil, 0, domain.NewInternalError(err)
	}
	return orderProducts, total, nil
}

// Create creates a new orderProduct
func (uc *orderProductUseCase) Create(ctx context.Context, i dto.CreateOrderProductInput) (*entity.OrderProduct, error) {
	orderProduct := i.ToEntity()

	if err := uc.gateway.Create(ctx, orderProduct); err != nil {
		return nil, domain.NewInternalError(err)
	}

	return orderProduct, nil
}

// Get returns a orderProduct by ID
func (uc *orderProductUseCase) Get(ctx context.Context, i dto.GetOrderProductInput) (*entity.OrderProduct, error) {
	orderProduct, err := uc.gateway.FindByID(ctx, i.OrderID, i.ProductID)
	if err != nil {
		return nil, domain.NewInternalError(err)
	}

	if orderProduct == nil {
		return nil, domain.NewNotFoundError(domain.ErrNotFound)
	}

	return orderProduct, nil
}

func (uc *orderProductUseCase) Update(ctx context.Context, i dto.UpdateOrderProductInput) (*entity.OrderProduct, error) {
	orderProduct, err := uc.gateway.FindByID(ctx, i.OrderID, i.ProductID)
	if err != nil {
		return nil, domain.NewInternalError(err)
	}

	if orderProduct == nil {
		return nil, domain.NewNotFoundError(domain.ErrNotFound)
	}

	order := orderProduct.Order
	product := orderProduct.Product
	orderProduct.Update(i.Quantity)

	if err := uc.gateway.Update(ctx, orderProduct); err != nil {
		return nil, domain.NewInternalError(err)
	}

	orderProduct.Order = order
	orderProduct.Product = product

	return orderProduct, nil
}

func (uc *orderProductUseCase) Delete(ctx context.Context, i dto.DeleteOrderProductInput) (*entity.OrderProduct, error) {
	order, err := uc.gateway.FindByID(ctx, i.OrderID, i.ProductID)
	if err != nil {
		return nil, domain.NewInternalError(err)
	}
	if order == nil {
		return nil, domain.NewNotFoundError(domain.ErrNotFound)
	}

	if err := uc.gateway.Delete(ctx, i.OrderID, i.ProductID); err != nil {
		return nil, domain.NewInternalError(err)
	}

	return order, nil
}
