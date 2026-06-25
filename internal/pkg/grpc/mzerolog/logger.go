package mzerolog

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/AlexBond702/catalog-service/internal/pkg/grpc/grpch"
)

type middleware struct {
	log zerolog.Logger
}

func NewMiddleware() grpch.Middleware {
	return &middleware{
		log: log.Logger,
	}
}

func (m *middleware) accessLogEvent(method string, start time.Time, err error) {
	level := zerolog.TraceLevel
	if err != nil {
		level = zerolog.ErrorLevel
	}

	m.log.WithLevel(level).
		Err(err).
		Int64("duration", time.Since(start).Milliseconds()).
		Str("grpc_method", method).
		Int("grpc_status", int(status.Convert(err).Code())).
		Send()
}

func (m *middleware) ForUnary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		start := time.Now()

		resp, err = handler(ctx, req)
		var logErr error
		if err != nil {
			logErr = err
		} else {
			logErr = ctx.Err()
		}
		m.accessLogEvent(info.FullMethod, start, logErr)
		return resp, err
	}
}

func (m *middleware) ForStream() grpc.StreamServerInterceptor {
	return func(
		srv any,
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) (err error) {
		ctx := ss.Context()
		start := time.Now()
		err = handler(srv, ss)
		var logErr error
		if err != nil {
			logErr = err
		} else {
			logErr = ctx.Err()
		}
		m.accessLogEvent(info.FullMethod, start, logErr)
		return err
	}
}
