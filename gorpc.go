package gorpc

import (
	"context"
	"fmt"
	"net"

	"github.com/kaptika/common/log"
	"github.com/kaptika/common/validator"
	"github.com/sethvargo/go-envconfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	GRPC *grpc.Server
}

type Config struct {
	Port int `env:"GORPC_PORT, default=4000"`
}

var config *Config
var server *Server

func init() {
	config = new(Config)

	// Load config from environment variables
	err := envconfig.Process(context.Background(), config)
	if err != nil {
		panic(err)
	}

	// Validate config
	err = validator.Validate(config)
	if err != nil {
		panic(err)
	}

	server = &Server{
		GRPC: grpc.NewServer(),
	}
	healthCheck := health.NewServer()
	grpc_health_v1.RegisterHealthServer(server.GRPC, healthCheck)
	reflection.Register(server.GRPC)
}

func GetGRPCServer() *grpc.Server {
	return server.GRPC
}

func Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port))
	if err != nil {
		panic(err)
	}

	log.Infof("gRPC server is listening at %s", lis.Addr())

	return server.GRPC.Serve(lis)
}

func Shutdown() error {
	server.GRPC.GracefulStop()

	return nil
}
