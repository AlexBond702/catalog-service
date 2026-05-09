package sproduct

import (
	"context"
	"database/sql"
	"errors"
	"github.com/AlexBond702/catalog-service/internal/app/entity"
	pcategory "github.com/AlexBond702/catalog-service/internal/app/repository/category"
	pproduct "github.com/AlexBond702/catalog-service/internal/app/repository/product"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type listProductSuite struct {
	suite.Suite
	service     *svc
	productRepo *pproduct.MockProduct
	ctx         context.Context
}
type getProductSuite struct {
	suite.Suite
	service     *svc
	productRepo *pproduct.MockProduct
	ctx         context.Context
}
type deleteProductSuite struct {
	suite.Suite
	service     *svc
	productRepo *pproduct.MockProduct
	ctx         context.Context
}
type updateProductSuite struct {
	suite.Suite
	service     *svc
	productRepo *pproduct.MockProduct
	ctx         context.Context
}
type createProductSuite struct {
	suite.Suite
	service      *svc
	productRepo  *pproduct.MockProduct
	categoryRepo *pcategory.MockCategory
	ctx          context.Context
}

func (l *listProductSuite) SetupTest() {
	l.ctx = context.Background()
	l.productRepo = new(pproduct.MockProduct)
	CategoryRepo := new(pcategory.MockCategory)
	l.service = NewService(l.productRepo, CategoryRepo).(*svc)
}

func (g *getProductSuite) SetupTest() {
	g.ctx = context.Background()
	g.productRepo = new(pproduct.MockProduct)
	CategoryRepo := new(pcategory.MockCategory)
	g.service = NewService(g.productRepo, CategoryRepo).(*svc)
}

func (d *deleteProductSuite) SetupTest() {
	d.ctx = context.Background()
	d.productRepo = new(pproduct.MockProduct)
	CategoryRepo := new(pcategory.MockCategory)
	d.service = NewService(d.productRepo, CategoryRepo).(*svc)
}

func (u *updateProductSuite) SetupTest() {
	u.ctx = context.Background()
	u.productRepo = new(pproduct.MockProduct)
	CategoryRepo := new(pcategory.MockCategory)
	u.service = NewService(u.productRepo, CategoryRepo).(*svc)
}
func (s *createProductSuite) SetupTest() {
	s.ctx = context.Background()
	s.productRepo = new(pproduct.MockProduct)
	s.categoryRepo = new(pcategory.MockCategory)
	s.service = NewService(s.productRepo, s.categoryRepo).(*svc)

}

func (l *listProductSuite) TearDownTest() {
	l.productRepo.AssertExpectations(l.T())
}

func (g *getProductSuite) TearDownTest() {
	g.productRepo.AssertExpectations(g.T())
}

func (d *deleteProductSuite) TearDownTest() {
	d.productRepo.AssertExpectations(d.T())
}

func (u *updateProductSuite) TearDownTest() {
	u.productRepo.AssertExpectations(u.T())
}
func (s *createProductSuite) TearDownTest() {
	s.productRepo.AssertExpectations(s.T())
	s.categoryRepo.AssertExpectations(s.T())
}

func TestListProductSuite(t *testing.T) {
	suite.Run(t, new(listProductSuite))
}

func TestGetProductSuite(t *testing.T) {
	suite.Run(t, new(getProductSuite))
}

func TestDeleteProductSuite(t *testing.T) {
	suite.Run(t, new(deleteProductSuite))
}

func TestUpdateProductSuite(t *testing.T) {
	suite.Run(t, new(updateProductSuite))
}
func TestCreateProductSuite(t *testing.T) {
	suite.Run(t, new(createProductSuite))
}

func (s *createProductSuite) TestCreateProduct() {
	type want struct {
		err    error
		result entity.Product
	}
	type args struct {
		request entity.RequestProductCreate
	}
	var Categoryguid = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	var Productguid = uuid.MustParse("22222222-2222-2222-2222-222222222222")

	testCases := []struct {
		name    string
		args    args
		want    want
		prepare func(args args)
	}{
		{
			name: "success",
			prepare: func(args args) {
				s.productRepo.EXPECT().
					List(s.ctx, &args.request.Name, &args.request.CategoryGUID).
					Return([]entity.Product{}, nil).
					Once()
				s.categoryRepo.EXPECT().
					GetByGUID(s.ctx, Categoryguid).
					Return(entity.Category{}, nil).
					Once()
				s.productRepo.EXPECT().
					Create(s.ctx, mock.MatchedBy(func(product entity.Product) bool {
						return product.Name == args.request.Name &&
							product.Description == args.request.Description &&
							product.Price == args.request.Price &&
							product.CategoryGUID == args.request.CategoryGUID
					})).
					Return(nil).
					Once()
			},
			args: args{
				entity.RequestProductCreate{
					Name:         "Test Product",
					Price:        100.50,
					CategoryGUID: Categoryguid,
				},
			},
			want: want{
				err: nil,
				result: entity.Product{
					GUID:         Productguid,
					Name:         "Test Product",
					Price:        100.50,
					CategoryGUID: Categoryguid,
				},
			},
		},
		{
			name: "already exist",
			prepare: func(args args) {
				s.productRepo.EXPECT().
					List(s.ctx, &args.request.Name, &args.request.CategoryGUID).
					Return([]entity.Product{{Name: args.request.Name}}, nil).
					Once()
			},
			args: args{
				request: entity.RequestProductCreate{
					Name:         "Test Product",
					Price:        100.50,
					CategoryGUID: Categoryguid,
				},
			},
			want: want{
				err:    entity.ErrAlreadyExists,
				result: entity.Product{},
			},
		},

		{
			name: "category not found",
			prepare: func(args args) {
				s.productRepo.EXPECT().
					List(s.ctx, &args.request.Name, &args.request.CategoryGUID).
					Return([]entity.Product{}, nil).
					Once()

				s.categoryRepo.EXPECT().GetByGUID(s.ctx, args.request.CategoryGUID).
					Return(entity.Category{}, entity.ErrNotFound).Once()

			},
			args: args{
				request: entity.RequestProductCreate{
					Name:         "Test Product",
					Price:        100.50,
					CategoryGUID: Categoryguid,
				},
			},
			want: want{
				err:    entity.ErrNotFound,
				result: entity.Product{},
			},
		},
	}
	for _, testCase := range testCases {
		s.Run(testCase.name, func() {
			testCase.prepare(testCase.args)
			result, err := s.service.Create(s.ctx, testCase.args.request)
			s.True(errors.Is(err, testCase.want.err), testCase.name)

			if testCase.want.err == nil {
				s.Equal(testCase.want.result.Name, result.Name)
				s.Equal(testCase.want.result.Description, result.Description)
				s.Equal(testCase.want.result.CategoryGUID, result.CategoryGUID)
				s.Equal(testCase.want.result.Price, result.Price)
				s.NotZero(result.CreatedAt)
				s.NotZero(result.UpdatedAt)
			} else {
				s.Empty(result.ID)
			}
		})
	}
}

func (u *updateProductSuite) TestUpdateProduct() {
	var CategoryGUID = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	var GUID = uuid.MustParse("33333333-3333-3333-3333-333333333333")
	var Description = "Description"
	var CreatedAt = time.Now().Add(-time.Hour)

	type want struct {
		expectedErr error
		result      entity.Product
	}
	type args struct {
		request entity.RequestProductUpdate
	}
	testCase := []struct {
		name    string
		want    want
		args    args
		prepare func(args args)
	}{
		{
			name: "success",
			want: want{
				expectedErr: nil,
				result: entity.Product{
					GUID:         GUID,
					Name:         "Test Update",
					Description:  &Description,
					Price:        1000.0,
					CategoryGUID: CategoryGUID,
				},
			},
			args: args{
				request: entity.RequestProductUpdate{
					Name:         "Test Update",
					Description:  &Description,
					Price:        1000.0,
					CategoryGUID: CategoryGUID,
				},
			},
			prepare: func(args args) {
				u.productRepo.EXPECT().
					GetByGUID(u.ctx, GUID).
					Return(entity.Product{
						GUID:      GUID,
						Name:      "Test Update",
						Price:     1000.0,
						CreatedAt: CreatedAt,
					}, nil).Once()
				u.productRepo.EXPECT().
					List(u.ctx, &args.request.Name, &args.request.CategoryGUID).
					Return([]entity.Product{}, nil).
					Once()
				u.productRepo.EXPECT().
					Update(u.ctx, mock.MatchedBy(func(product entity.Product) bool {
						return product.Name == args.request.Name &&
							product.Description == args.request.Description &&
							product.Price == args.request.Price &&
							product.CategoryGUID == args.request.CategoryGUID
					})).Return(nil).Once()
			},
		},
		{
			name: "product not found",
			want: want{
				expectedErr: entity.ErrNotFound,
				result:      entity.Product{},
			},
			args: args{
				request: entity.RequestProductUpdate{
					Name:         "Test Update",
					Description:  &Description,
					Price:        1000.0,
					CategoryGUID: CategoryGUID,
				},
			},
			prepare: func(args args) {
				u.productRepo.EXPECT().GetByGUID(u.ctx, GUID).
					Return(entity.Product{}, sql.ErrNoRows).
					Once()
			},
		},
		{
			name: "category already exist",
			want: want{
				expectedErr: entity.ErrAlreadyExists,
				result:      entity.Product{},
			},
			args: args{
				request: entity.RequestProductUpdate{
					Name:         "Test Update",
					Description:  &Description,
					Price:        1000.0,
					CategoryGUID: CategoryGUID,
				},
			},
			prepare: func(args args) {
				u.productRepo.EXPECT().
					GetByGUID(u.ctx, GUID).
					Return(entity.Product{
						CategoryGUID: args.request.CategoryGUID,
					}, nil).Once()
				u.productRepo.EXPECT().
					List(u.ctx, &args.request.Name, &args.request.CategoryGUID).
					Return([]entity.Product{
						{
							Name:         args.request.Name,
							CategoryGUID: args.request.CategoryGUID,
						},
					}, nil).Once()
			},
		},
	}

	for _, test := range testCase {
		u.Run(test.name, func() {
			test.prepare(test.args)
			result, err := u.service.Update(u.ctx, GUID, test.args.request)
			u.True(errors.Is(err, test.want.expectedErr), test.name)

			if test.want.expectedErr == nil {
				u.Equal(test.args.request.Name, result.Name)
				u.Equal(test.args.request.Price, result.Price)
				u.Equal(test.args.request.Description, result.Description)
				u.Equal(test.args.request.CategoryGUID, result.CategoryGUID)
				u.NotZero(result.UpdatedAt)
				u.NotZero(result.CreatedAt)
			} else {
				u.Empty(result.ID)
			}
		})
	}
}

func (d *deleteProductSuite) TestDeleteProduct() {
	var GUID = uuid.MustParse("33333333-3333-3333-3333-333333333333")

	type want struct {
		expectedErr error
		result      entity.Product
	}
	type args struct {
		request entity.Product
	}
	testCase := []struct {
		name    string
		want    want
		args    args
		prepare func(args args)
	}{
		{
			name: "success",
			want: want{
				expectedErr: nil,
				result:      entity.Product{},
			},
			args: args{
				entity.Product{
					GUID: GUID,
				},
			},
			prepare: func(args args) {
				d.productRepo.EXPECT().
					GetByGUID(d.ctx, GUID).
					Return(entity.Product{
						GUID: GUID,
					}, nil).Once()
				d.productRepo.EXPECT().
					Delete(d.ctx, GUID).
					Return(nil).Once()
			},
		},
		{
			name: "product not found",
			want: want{
				expectedErr: entity.ErrNotFound,
				result:      entity.Product{},
			},
			args: args{
				request: entity.Product{
					GUID: GUID,
				},
			},
			prepare: func(args args) {
				d.productRepo.EXPECT().
					GetByGUID(d.ctx, GUID).
					Return(entity.Product{}, sql.ErrNoRows).
					Once()
			},
		},
	}
	for _, test := range testCase {
		d.Run(test.name, func() {
			test.prepare(test.args)
			err := d.service.Delete(d.ctx, GUID)
			d.True(errors.Is(err, test.want.expectedErr), test.name)

		})
	}
}

func (g *getProductSuite) TestGetProduct() {
	var GUID = uuid.MustParse("33333333-3333-3333-3333-333333333333")

	type want struct {
		expectedErr error
		result      entity.Product
	}
	type args struct {
		request entity.Product
	}
	testCase := []struct {
		name    string
		want    want
		args    args
		prepare func(args args)
	}{
		{
			name: "success",
			want: want{
				expectedErr: nil,
				result: entity.Product{
					GUID: GUID,
				},
			},
			args: args{
				request: entity.Product{
					GUID: GUID,
				},
			},
			prepare: func(args args) {
				g.productRepo.EXPECT().
					GetByGUID(g.ctx, args.request.GUID).
					Return(entity.Product{
						GUID: GUID,
					}, nil).Once()
			},
		},
		{
			name: "product not found",
			want: want{
				expectedErr: entity.ErrNotFound,
				result:      entity.Product{},
			},
			args: args{
				request: entity.Product{
					GUID: GUID,
				},
			},
			prepare: func(args args) {
				g.productRepo.EXPECT().
					GetByGUID(g.ctx, args.request.GUID).
					Return(entity.Product{}, sql.ErrNoRows).Once()
			},
		},
	}
	for _, test := range testCase {
		g.Run(test.name, func() {
			test.prepare(test.args)
			result, err := g.service.GetByGUID(g.ctx, GUID)
			g.True(errors.Is(err, test.want.expectedErr), test.name)

			if test.want.expectedErr == nil {
				g.Equal(test.want.result.Name, result.Name)
				g.Equal(test.want.result.Price, result.Price)
				g.Equal(test.want.result.GUID, result.GUID)
				g.Equal(test.want.result.CategoryGUID, result.CategoryGUID)
			} else {
				g.Empty(result.ID)
			}
		})
	}
}
func (l *listProductSuite) TestListProduct() {
	var (
		name     *string    = nil
		category *uuid.UUID = nil
		ListErr             = errors.New("list error")
	)
	type want struct {
		expectedErr error
		result      []entity.Product
	}

	testCase := []struct {
		name    string
		want    want
		prepare func()
	}{
		{
			name: "success",
			want: want{
				expectedErr: nil,
				result: []entity.Product{
					{
						Name:  "Test Product 1",
						Price: 100.0,
					},
					{
						Name:  "Test Product 2",
						Price: 150.0,
					},
				},
			},
			prepare: func() {
				l.productRepo.EXPECT().
					List(l.ctx, name, category).
					Return([]entity.Product{
						{
							Name:  "Test Product 1",
							Price: 100.0,
						},
						{
							Name:  "Test Product 2",
							Price: 150.0,
						},
					}, nil).Once()
			},
		},
		{
			name: "list error",
			want: want{
				expectedErr: ListErr,
				result:      []entity.Product{},
			},
			prepare: func() {
				l.productRepo.EXPECT().
					List(l.ctx, name, category).
					Return([]entity.Product{}, ListErr).
					Once()
			},
		},
	}
	for _, test := range testCase {
		l.Run(test.name, func() {
			test.prepare()
			result, err := l.service.List(l.ctx)
			if test.want.expectedErr != nil {
				l.ErrorIs(err, test.want.expectedErr)
				return
			}
			l.NoError(err)
			for i := range test.want.result {
				l.Equal(test.want.result[i].Name, result[i].Name)
				l.Equal(test.want.result[i].Price, result[i].Price)
				l.Equal(test.want.result[i].CategoryGUID, result[i].CategoryGUID)
			}
		})
	}
}
