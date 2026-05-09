package entity

import "errors"

var (
	ErrNotFound            = errors.New("not found")
	ErrAlreadyExists       = errors.New("already exists")
	ErrCategoryHasProducts = errors.New("category has linked products")
	ErrIncorrectParameters = errors.New("incorrect parameters")
	ErrProductDuplicate    = errors.New("product duplicate")
	ErrCategoryDuplicate   = errors.New("category duplicate")
	ErrInvalidPrice        = errors.New("price must be positive")
	ErrInvalidName         = errors.New("product name cannot be empty")
)
