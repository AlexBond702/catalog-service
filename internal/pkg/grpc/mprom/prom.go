package mprom

import (
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"

	"github.com/AlexBond702/catalog-service/internal/pkg/grpc/grpch"
)

type middleware struct{}

func New() grpch.Middleware {
	return new(middleware)
}

func (m *middleware) ForUnary() grpc.UnaryServerInterceptor {
	return grpc_prometheus.UnaryServerInterceptor
}

func (m *middleware) ForStream() grpc.StreamServerInterceptor {
	return grpc_prometheus.StreamServerInterceptor
}

func EnableHistogram() {
	grpc_prometheus.EnableHandlingTimeHistogram()
}
