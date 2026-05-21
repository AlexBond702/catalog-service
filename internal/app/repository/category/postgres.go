package pcategory

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

func NewRepoFromPostgres(ctx context.Context, client *rcpostgres.Client) repository.Category {
	return &repoPg{_DB: client}
}

func (r *repoPg) Create(ctx context.Context, category entity.Category) error {
	_, err := r._DB.NewInsert().Model(&category).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *repoPg) GetByGUID(ctx context.Context, guid uuid.UUID) (entity.Category, error) {
	var cat entity.Category
	err := r._DB.NewSelect().Model(&cat).Where("guid = ?", guid).Scan(ctx)
	if err != nil {
		return entity.Category{}, util.ReplaceErr1(err, sql.ErrNoRows, entity.ErrNotFound)
	}
	return cat, nil
}

func (r *repoPg) Update(ctx context.Context, category entity.Category) error {
	err := r._DB.NewUpdate().Model(&category).WherePK().ExcludeColumn("id", "created_at").Returning("*").Scan(ctx)
	return err
}

func (r *repoPg) Delete(ctx context.Context, guid uuid.UUID) error {
	var deleted entity.Category
	_, err := r._DB.NewDelete().Model((*entity.Category)(nil)).Where("guid = ?", guid).Exec(ctx, &deleted)
	return err
}

func (r *repoPg) List(ctx context.Context, name *string) ([]entity.Category, error) {
	var categories []entity.Category
	query := r._DB.NewSelect().Model(&categories)
	if name != nil {
		query = query.Where("name = ?", *name)
	}
	err := query.Scan(ctx)
	if err != nil {
		return nil, err
	}
	return categories, nil
}
