package router

import (
	"google.golang.org/grpc"

	"github.com/Te8va/shortURL/internal/app/config"
	pb "github.com/Te8va/shortURL/internal/app/grpc/gen"
	"github.com/Te8va/shortURL/internal/app/grpchandler"
	"github.com/Te8va/shortURL/internal/app/service"
)

func NewGRPCRouter(
	cfg *config.Config,
	saver service.URLSaverServ,
	getter service.URLGetterServ,
	deleter service.URLDeleteServ,
	pinger service.PingerServ,
	stat service.URLStatsServ,
) *grpc.Server {

	server := grpc.NewServer()

	handler := grpchandler.NewShortURLHandler(
		saver,
		getter,
		deleter,
		pinger,
		stat,
		cfg,
	)

	pb.RegisterShortURLServiceServer(server, handler)
	return server
}
