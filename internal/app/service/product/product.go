package sproduct

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/AlexBond702/catalog-service/internal/app/entity"
	"github.com/AlexBond702/catalog-service/internal/app/repository"
	"github.com/AlexBond702/catalog-service/internal/app/service"
)

type svc struct {
	repoCategory repository.Category
	repoProduct  repository.Product
}

func NewService(repoProduct repository.Product, repoCategory repository.Category) service.Product {
	return &svc{repoProduct: repoProduct, repoCategory: repoCategory}
}

func (s *svc) Create(ctx context.Context, req entity.RequestProductCreate) (entity.Product, error) {
	products, err := s.repoProduct.List(ctx, &req.Name, &req.CategoryGUID)
	if err != nil {
		return entity.Product{}, err
	}
	if len(products) > 0 {
		return entity.Product{}, entity.ErrAlreadyExists
	}
	_, err = s.repoCategory.GetByGUID(ctx, req.CategoryGUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Product{}, entity.ErrNotFound
		}
		return entity.Product{}, err
	}
	now := time.Now()
	newProduct := entity.Product{
		GUID:         uuid.Must(uuid.NewRandom()),
		Name:         req.Name,
		Description:  req.Description,
		Price:        req.Price,
		CategoryGUID: req.CategoryGUID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := s.repoProduct.Create(ctx, newProduct); err != nil {
		return entity.Product{}, err
	}
	return newProduct, nil
}

func (s *svc) GetByGUID(ctx context.Context, guid uuid.UUID) (entity.Product, error) {
	gettingProd, err := s.repoProduct.GetByGUID(ctx, guid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Product{}, entity.ErrNotFound
		}
		return entity.Product{}, err
	}
	return gettingProd, nil
}

func (s *svc) Update(ctx context.Context, guid uuid.UUID, req entity.RequestProductUpdate) (entity.Product, error) {
	oldProduct, err := s.repoProduct.GetByGUID(ctx, guid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Product{}, entity.ErrNotFound
		}
		return entity.Product{}, err
	}
	updatingProduct, err := s.repoProduct.List(ctx, &req.Name, &req.CategoryGUID)
	if err != nil {
		return entity.Product{}, err
	}
	if len(updatingProduct) > 0 {
		return entity.Product{}, entity.ErrAlreadyExists
	}
	now := time.Now()
	oldProduct.Name = req.Name
	oldProduct.CategoryGUID = req.CategoryGUID
	oldProduct.Price = req.Price
	oldProduct.Description = req.Description
	oldProduct.UpdatedAt = now

	if err := s.repoProduct.Update(ctx, oldProduct); err != nil {
		return entity.Product{}, err
	}
	return oldProduct, nil
}

func (s *svc) Delete(ctx context.Context, guid uuid.UUID) error {
	product, err := s.repoProduct.GetByGUID(ctx, guid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.ErrNotFound
		}
		return err
	}
	if err := s.repoProduct.Delete(ctx, product.GUID); err != nil {
		return err
	}
	return nil
}

func (s *svc) List(ctx context.Context) ([]entity.Product, error) {
	products, err := s.repoProduct.List(ctx, nil, nil)
	if err != nil {
		return nil, err
	}
	return products, nil
}
