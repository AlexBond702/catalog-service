package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/AlexBond702/catalog-service/internal/app/entity"
)

type Transactional interface {
	InsideTx(ctx context.Context, cb func(ctx context.Context) error) error
}

type Category interface {
	Transactional
	Create(ctx context.Context, category entity.Category) error
	GetByGUID(ctx context.Context, guid uuid.UUID) (entity.Category, error)
	Update(ctx context.Context, category entity.Category) error
	Delete(ctx context.Context, guid uuid.UUID) error
	List(ctx context.Context, name *string) ([]entity.Category, error)
}

type Product interface {
	Transactional
	Create(ctx context.Context, product entity.Product) error
	GetByGUID(ctx context.Context, guid uuid.UUID) (entity.Product, error)
	Update(ctx context.Context, product entity.Product) error
	Delete(ctx context.Context, guid uuid.UUID) error
	List(ctx context.Context, name *string, categoryGUID *uuid.UUID) ([]entity.Product, error)
}
