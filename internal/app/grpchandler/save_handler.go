package grpchandler

import (
	"context"
	"errors"
	"net/url"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Te8va/shortURL/internal/app/domain"
	appErrors "github.com/Te8va/shortURL/internal/app/errors"
	gen "github.com/Te8va/shortURL/internal/app/grpc/gen"
)

// SaveURL processes gRPC requests to save URL.
func (h *ShortURLHandler) SaveURL(ctx context.Context, req *gen.SaveURLRequest) (*gen.SaveURLResponse, error) {
	userID, _ := ctx.Value(domain.UserIDKey).(int)

	if req.Url == "" {
		return nil, status.Error(codes.InvalidArgument, "url is required")
	}

	_, err := url.ParseRequestURI(req.Url)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid URL format")
	}

	shortURL, err := h.saver.Save(ctx, userID, req.Url)
	if err != nil {
		if errors.Is(err, appErrors.ErrURLExists) {
			return &gen.SaveURLResponse{ShortUrl: shortURL}, nil
		}
		return nil, status.Error(codes.Internal, "failed to save URL")
	}

	return &gen.SaveURLResponse{ShortUrl: shortURL}, nil
}

// SaveBatch processes gRPC batch URL saving requests.
func (h *ShortURLHandler) SaveBatch(ctx context.Context, req *gen.SaveBatchRequest) (*gen.SaveBatchResponse, error) {
	userID, _ := ctx.Value(domain.UserIDKey).(int)

	if len(req.Items) == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty batch request")
	}

	urlMap := make(map[string]string)
	for _, item := range req.Items {
		urlMap[item.CorrelationId] = item.OriginalUrl
	}

	result, err := h.saver.SaveBatch(ctx, userID, urlMap)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to save batch URLs")
	}

	var responseItems []*gen.BatchResponseItem
	for correlationID, shortURL := range result {
		responseItems = append(responseItems, &gen.BatchResponseItem{
			CorrelationId: correlationID,
			ShortUrl:      shortURL,
		})
	}

	return &gen.SaveBatchResponse{
		Items: responseItems,
	}, nil
}
