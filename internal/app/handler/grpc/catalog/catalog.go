package catalog

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/AlexBond702/catalog-service/gen/proto/catalog/v1"
	"github.com/AlexBond702/catalog-service/internal/app/entity"
	"github.com/AlexBond702/catalog-service/internal/app/service"
)

type Handler struct {
	pb.UnimplementedCatalogServiceServer
	productService service.Product
}

func NewHandler(productService service.Product) *Handler {
	return &Handler{
		productService: productService,
	}
}

func (h *Handler) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	guid, err := uuid.Parse(req.Guid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	product, err := h.productService.GetByGUID(ctx, guid)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	pbProduct := convertProductToProto(product)
	return &pb.GetProductResponse{
		Product: pbProduct,
	}, nil
}

func (h *Handler) GetProducts(ctx context.Context, req *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {
	productsProto := make([]*pb.Product, 0)
	for _, guid := range req.Guids {
		guidUUID, err := uuid.Parse(guid)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		product, err := h.productService.GetByGUID(ctx, guidUUID)
		if err != nil {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		productProto := convertProductToProto(product)
		productsProto = append(productsProto, productProto)
	}
	return convertProductsToProto(productsProto), nil
}

func (h *Handler) CheckProductExists(ctx context.Context, req *pb.CheckProductExistsRequest) (*pb.CheckProductExistsResponse, error) {
	guidStr, err := uuid.Parse(req.Guid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	product, err := h.productService.GetByGUID(ctx, guidStr)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	check := convertCheckProduct(product)
	return &pb.CheckProductExistsResponse{
		Product: check,
	}, nil
}

func convertProductToProto(product entity.Product) *pb.Product {
	return &pb.Product{
		Id:          product.ID,
		Guid:        product.GUID.String(),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CreatedAt:   timestamppb.New(product.CreatedAt),
		UpdatedAt:   timestamppb.New(product.UpdatedAt),
	}
}

func convertProductsToProto(products []*pb.Product) *pb.GetProductsResponse {
	return &pb.GetProductsResponse{
		Products: products,
	}
}

func convertCheckProduct(product entity.Product) *pb.Product {
	return &pb.Product{
		Price: product.Price,
	}
}
