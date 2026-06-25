package mcommon

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/AlexBond702/catalog-service/internal/pkg/grpc/grpch"
)

type recoveryMiddleware struct{}

func NewRecovery() grpch.Middleware {
	return new(recoveryMiddleware)
}

func (r *recoveryMiddleware) ForUnary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		panicked := true
		defer func() {
			if rec := recover(); rec != nil || panicked {
				err = status.Errorf(codes.Internal, "%s", fmt.Sprintf("recovery panic: %v - %t", rec, panicked))
			}
		}()
		resp, err = handler(ctx, req)
		panicked = false
		return resp, err
	}
}

func (r *recoveryMiddleware) ForStream() grpc.StreamServerInterceptor {
	return func(
		srv any,
		cc grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) (err error) {
		panicked := true
		defer func() {
			if rec := recover(); rec != nil || panicked {
				err = status.Errorf(codes.Internal, "%s", fmt.Sprintf("recovery panic: %v - %t", rec, panicked))
			}
		}()
		err = handler(srv, cc)
		panicked = false
		return err
	}
}
