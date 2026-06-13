package pproduct

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	"github.com/AlexBond702/catalog-service/internal/app/entity"
	"github.com/AlexBond702/catalog-service/internal/app/repository"
	rcpostgres "github.com/AlexBond702/catalog-service/internal/app/repository/conn/postgres"
	"github.com/AlexBond702/catalog-service/internal/app/util"
)

type (
	repoPg struct {
		*_DB
	}
	_DB = rcpostgres.Client
)

func NewRepoFromPostgres(ctx context.Context, client *rcpostgres.Client) repository.Product {
	return &repoPg{_DB: client}
}

func (r *repoPg) Create(ctx context.Context, product entity.Product) error {
	_, err := r._DB.NewInsert().Model(&product).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *repoPg) GetByGUID(ctx context.Context, guid uuid.UUID) (entity.Product, error) {
	var prod entity.Product
	err := r._DB.NewSelect().Model(&prod).Where("guid = ?", guid).Scan(ctx)
	if err != nil {
		return entity.Product{}, util.ReplaceErr1(err, sql.ErrNoRows, entity.ErrNotFound)
	}
	return prod, nil
}

func (r *repoPg) GetById(ctx context.Context, id int64) (entity.Product, error) {
	var product entity.Product
	if err := r._DB.NewSelect().Model(&product).Where("id=?", id).Scan(ctx); err != nil {
		return entity.Product{}, util.ReplaceErr1(err, sql.ErrNoRows, entity.ErrNotFound)
	}
	return product, nil
}

func (r *repoPg) Update(ctx context.Context, product entity.Product) error {
	err := r._DB.NewUpdate().Model(&product).WherePK().ExcludeColumn("id", "created_at").Returning("*").Scan(ctx)
	return err
}

func (r *repoPg) Delete(ctx context.Context, guid uuid.UUID) error {
	var deleted entity.Product
	err := r._DB.NewDelete().Model((*entity.Product)(nil)).Where("guid =?", guid).Returning("*").Scan(ctx, &deleted)
	return err
}

func (r *repoPg) List(ctx context.Context, name *string, categoryGUID *uuid.UUID) ([]entity.Product, error) {
	var products []entity.Product
	query := r._DB.NewSelect().Model(&products)

	if name != nil {
		query = query.Where("name = ?", *name)
	}

	if categoryGUID != nil {
		query = query.Where("category_guid = ?", *categoryGUID)
	}
	err := query.Scan(ctx)
	if err != nil {
		return nil, err
	}
	return products, nil
}
