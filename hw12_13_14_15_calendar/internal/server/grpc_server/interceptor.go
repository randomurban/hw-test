package grpc_server

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func (s *Server) RequestLogInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	start := time.Now()
	resp, err = handler(ctx, req)
	duration := time.Since(start)
	s.logger.Info("grpc:", "method", info.FullMethod, "status", status.Code(err), "duration", duration, "err", err)
	return resp, err
}
