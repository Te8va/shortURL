package grpchandler

import (
	"github.com/Te8va/shortURL/internal/app/config"
	gen "github.com/Te8va/shortURL/internal/app/grpc/gen"
	"github.com/Te8va/shortURL/internal/app/service"
)

// ShortURLHandler handles all gRPC requests for the URL shortener service
type ShortURLHandler struct {
	gen.UnimplementedShortURLServiceServer
	saver   service.URLSaverServ
	getter  service.URLGetterServ
	deleter service.URLDeleteServ
	pinger  service.PingerServ
	stat    service.URLStatsServ
	cfg     *config.Config
}

// NewShortURLHandler creates a new instance of ShortURLHandler with all required dependencies
func NewShortURLHandler(
	saver service.URLSaverServ,
	getter service.URLGetterServ,
	deleter service.URLDeleteServ,
	pinger service.PingerServ,
	stat service.URLStatsServ,
	cfg *config.Config,
) *ShortURLHandler {
	return &ShortURLHandler{
		saver:   saver,
		getter:  getter,
		deleter: deleter,
		pinger:  pinger,
		stat:    stat,
		cfg:     cfg,
	}
}
