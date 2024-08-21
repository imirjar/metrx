package grpc

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/imirjar/metrx/internal/api"
	"github.com/imirjar/metrx/internal/models"
	"google.golang.org/grpc"
)

type Service interface {
	UpdateMetrics(ctx context.Context, metrics []models.Metrics) error
	UpdateMetric(ctx context.Context, metric models.Metrics) (models.Metrics, error)
	ViewMetric(ctx context.Context, metric models.Metrics) (models.Metrics, error)
	ViewMetrics(ctx context.Context) (map[string][]models.Metrics, error)
}

type GRPCServer struct {
	Service Service
	Server  *grpc.Server
}

func NewGRPCServer() *GRPCServer {
	gtw := GRPCServer{}
	return &gtw
}

func (gs *GRPCServer) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", "3200"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	// api.RegisterGoMetricsServer(grpcServer, api.UnimplementedGoMetricsServer{})
	return grpcServer.Serve(lis)
}

func (gs *GRPCServer) Stop() {
	gs.Server.Stop()
}

func (gs *GRPCServer) Update(ctx context.Context, r *api.Request) (*api.Response, error) {
	var w api.Response

	for _, m := range r.Metrics {
		log.Print(m)
	}

	return &w, nil
}
