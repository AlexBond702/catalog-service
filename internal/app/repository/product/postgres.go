package pproduct

import (
	"context"
	"database/sql"
	"fmt"

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

func NewRepoFromPostgres(client *rcpostgres.Client) repository.Product {
	return &repoPg{_DB: client}
}

func (r *repoPg) Create(ctx context.Context, product entity.Product) error {
	_, err := r.GetRawBunDb().NewInsert().Model(&product).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to create: %w", err)
	}
	return nil
}

func (r *repoPg) GetByGUID(ctx context.Context, guid uuid.UUID) (entity.Product, error) {
	var prod entity.Product
	err := r.GetRawBunDb().NewSelect().Model(&prod).Where("guid = ?", guid).Scan(ctx)
	if err != nil {
		return entity.Product{}, util.ReplaceErr1(err, sql.ErrNoRows, entity.ErrNotFound)
	}
	return prod, nil
}

func (r *repoPg) Update(ctx context.Context, product entity.Product) error {
	result, err := r.GetRawBunDb().NewUpdate().Model(&product).WherePK().ExcludeColumn("id", "created_at").Exec(ctx)
	if err != nil {
		return rcpostgres.UpdateErr(result, err)
	}
	return nil
}

func (r *repoPg) Delete(ctx context.Context, guid uuid.UUID) error {
	_, err := r.GetRawBunDb().NewDelete().Model((*entity.Product)(nil)).Where("guid =?", guid).Exec(ctx)
	if err != nil {
		return rcpostgres.DeleteErr(err)
	}
	return nil
}

func (r *repoPg) List(ctx context.Context, name *string, categoryGUID *uuid.UUID) ([]entity.Product, error) {
	var products []entity.Product
	query := r.GetRawBunDb().NewSelect().Model(&products)

	if name != nil {
		query = query.Where("name = ?", *name)
	}

	if categoryGUID != nil {
		query = query.Where("guid = ?", *categoryGUID)
	}
	err := query.Scan(ctx)
	if err != nil {
		return products, err
	}
	return products, nil
}
