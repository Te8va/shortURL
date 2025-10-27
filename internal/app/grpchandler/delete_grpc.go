package grpchandler

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Te8va/shortURL/internal/app/domain"
	gen "github.com/Te8va/shortURL/internal/app/grpc/gen"
)

// DeleteUserURLs processes gRPC requests to delete user URLs.
func (h *ShortURLHandler) DeleteUserURLs(ctx context.Context, req *gen.DeleteUserURLsRequest) (*gen.DeleteUserURLsResponse, error) {
	userID, ok := ctx.Value(domain.UserIDKey).(int)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user authentication required")
	}

	if req.Ids == nil || len(req.Ids) == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty list of URL IDs")
	}

	var fullURLs []string
	for _, id := range req.Ids {
		fullURLs = append(fullURLs, fmt.Sprintf("%s/%s", h.cfg.BaseURL, id))
	}

	go func(urls []string, uid int) {
		bgCtx := context.Background()
		if err := h.deleter.DeleteUserURLs(bgCtx, urls, uid); err != nil {
			fmt.Printf("Failed to delete URLs for user %d: %v\n", uid, err)
		}
	}(fullURLs, userID)

	return &gen.DeleteUserURLsResponse{}, nil
}
