package entity

import "errors"

var (
	ErrNotFound            = errors.New("not found")
	ErrAlreadyExists       = errors.New("already exists")
	ErrCategoryHasProducts = errors.New("category has linked products")
	ErrIncorrectParameters = errors.New("incorrect parameters")
	ErrProductDuplicate    = errors.New("product duplicate")
	ErrCategoryDuplicate   = errors.New("category duplicate")
)
