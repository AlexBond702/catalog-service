package scategory

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

func NewService(repoCategory repository.Category, repoProduct repository.Product) service.Category {
	return &svc{repoCategory: repoCategory, repoProduct: repoProduct}
}

func (s *svc) Create(ctx context.Context, req entity.RequestCategoryCreate) (entity.Category, error) {
	existing, err := s.repoCategory.List(ctx, &req.Name)
	if err != nil {
		return entity.Category{}, err
	}
	if len(existing) > 0 {
		return entity.Category{}, entity.ErrAlreadyExists
	}
	now := time.Now()
	category := entity.Category{
		GUID:      uuid.Must(uuid.NewRandom()),
		Name:      req.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.repoCategory.Create(ctx, category); err != nil {
		return entity.Category{}, err
	}
	return category, nil
}

func (s *svc) GetByGUID(ctx context.Context, guid uuid.UUID) (entity.Category, error) {
	getting, err := s.repoCategory.GetByGUID(ctx, guid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Category{}, entity.ErrNotFound
		}
		return entity.Category{}, err
	}
	return getting, nil
}

func (s *svc) Update(ctx context.Context, guid uuid.UUID, req entity.RequestCategoryUpdate) (entity.Category, error) {
	getting, err := s.repoCategory.GetByGUID(ctx, guid)
	if err != nil {
		return entity.Category{}, entity.ErrNotFound
	}
	existing, err := s.repoCategory.List(ctx, &req.Name)
	if err != nil {
		return entity.Category{}, err
	}
	if len(existing) > 0 {
		return entity.Category{}, entity.ErrAlreadyExists
	}
	now := time.Now()
	getting.Name = req.Name
	getting.UpdatedAt = now
	if err := s.repoCategory.Update(ctx, getting); err != nil {
		return entity.Category{}, err
	}
	return getting, nil
}

func (s *svc) Delete(ctx context.Context, guid uuid.UUID) error {
	_, err := s.repoCategory.GetByGUID(ctx, guid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.ErrNotFound
		}
		return err
	}
	products, err := s.repoProduct.List(ctx, nil, &guid)
	if err != nil {
		return err
	}
	if len(products) > 0 {
		return entity.ErrCategoryHasProducts
	}
	if err := s.repoCategory.Delete(ctx, guid); err != nil {
		return err
	}
	return nil
}

func (s *svc) List(ctx context.Context) ([]entity.Category, error) {
	categories, err := s.repoCategory.List(ctx, nil)
	if err != nil {
		return nil, err
	}
	return categories, nil
}
