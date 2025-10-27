package grpchandler

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Te8va/shortURL/internal/app/domain"
	gen "github.com/Te8va/shortURL/internal/app/grpc/gen"
)

// GetURL processes gRPC request to get original URL by short ID.
func (h *ShortURLHandler) GetURL(ctx context.Context, req *gen.GetURLRequest) (*gen.GetURLResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	fullID := fmt.Sprintf("%s/%s", h.cfg.BaseURL, req.Id)
	originalURL, exists, isDeleted := h.getter.Get(ctx, fullID)

	if !exists {
		return nil, status.Error(codes.NotFound, "URL not found")
	}

	if isDeleted {
		return nil, status.Error(codes.NotFound, "URL has been deleted")
	}

	return &gen.GetURLResponse{
		OriginalUrl: originalURL,
	}, nil
}

// GetUserURLs processes gRPC request to retrieve all URLs created by user.
func (h *ShortURLHandler) GetUserURLs(ctx context.Context, req *gen.GetUserURLsRequest) (*gen.GetUserURLsResponse, error) {
	userID, ok := ctx.Value(domain.UserIDKey).(int)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user authentication required")
	}

	urls, err := h.getter.GetUserURLs(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get user URLs")
	}

	if urls == nil {
		return &gen.GetUserURLsResponse{Urls: []*gen.UserURLItem{}}, nil
	}

	var userURLs []*gen.UserURLItem
	for _, url := range urls {
		for shortURL, originalURL := range url {
			userURLs = append(userURLs, &gen.UserURLItem{
				ShortUrl:    shortURL,
				OriginalUrl: originalURL,
			})
		}
	}

	return &gen.GetUserURLsResponse{
		Urls: userURLs,
	}, nil
}
