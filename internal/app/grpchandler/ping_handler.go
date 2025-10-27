package grpchandler

import (
	"context"

	gen "github.com/Te8va/shortURL/internal/app/grpc/gen"
)

// Ping processes gRPC request to health check of database connection.
func (h *ShortURLHandler) Ping(ctx context.Context, req *gen.PingRequest) (*gen.PingResponse, error) {
	err := h.pinger.PingPg(ctx)
	if err != nil {
		return &gen.PingResponse{Success: false}, nil
	}

	return &gen.PingResponse{Success: true}, nil
}
