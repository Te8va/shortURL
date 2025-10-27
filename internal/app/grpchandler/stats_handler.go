package grpchandler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gen "github.com/Te8va/shortURL/internal/app/grpc/gen"
)

// GetStats processes gRPC requests to retrieve service statistics
func (h *ShortURLHandler) GetGRPCStats(ctx context.Context, req *gen.GetStatsRequest) (*gen.GetStatsResponse, error) {
	urlsCount, usersCount, err := h.stat.GetStats(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get stats")
	}

	return &gen.GetStatsResponse{
		Urls:  int32(urlsCount),
		Users: int32(usersCount),
	}, nil
}
