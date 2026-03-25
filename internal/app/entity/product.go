package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Product struct {
	bun.BaseModel `bun:"table:product"`
	ID            int64     `bun:"id,autoincrement"`
	GUID          uuid.UUID `bun:"guid,pk"`
	Name          string    `bun:"name,notnull"`
	Description   *string   `bun:"description"`
	Price         float64   `bun:"price,notnull"`
	CategoryGUID  uuid.UUID `bun:"category_guid,notnull"`
	CreatedAt     time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt     time.Time `bun:"updated_at,notnull,default:current_timestamp"`
}

type RequestProductCreate struct {
	Name         string    `json:"name" binding:"required,min=2,max=255"`
	Description  *string   `json:"description" binding:"omitempty,max=255"`
	Price        float64   `json:"price" binding:"required,numeric,ne=0,gt=0"`
	CategoryGUID uuid.UUID `json:"category_guid" binding:"required,uuid"`
}

type RequestProductUpdate struct {
	Name         string    `json:"name" binding:"required,min=2,max=255"`
	Description  *string   `json:"description" binding:"omitempty,max=255"`
	Price        float64   `json:"price" binding:"required,numeric,ne=0,gt=0"`
	CategoryGUID uuid.UUID `json:"category_guid" binding:"required,uuid"`
}

type ResponseProduct struct {
	GUID         uuid.UUID `json:"guid"`
	Name         string    `json:"name"`
	Description  *string   `json:"description"`
	Price        float64   `json:"price"`
	CategoryGUID uuid.UUID `json:"category_guid"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
