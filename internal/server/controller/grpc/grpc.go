package grpc

import (
	"context"
	"errors"
	"io"
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
	api.UnimplementedGoMetricsServer
	Service Service
}

func New() *GRPCServer {
	return &GRPCServer{}
}

func (gs *GRPCServer) Start(addr string) error {
	var opts []grpc.ServerOption
	s := grpc.NewServer(opts...)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	api.RegisterGoMetricsServer(s, gs)

	log.Printf("RUN GRPC on: %s", addr)
	return s.Serve(lis)
}

func (gs *GRPCServer) BatchUpdate(stream api.GoMetrics_BatchUpdateServer) error {
	metrics := []models.Metrics{}

	for {
		metric, err := stream.Recv()

		if err != nil {
			if err == io.EOF {
				err = gs.Service.UpdateMetrics(context.Background(), metrics)
				if err != nil {
					log.Print(err)
				}
				return stream.SendAndClose(&api.Response{})
			}
			return err
		}

		switch metric.Type {
		case "gauge":
			value := &metric.Value
			var m = models.Metrics{
				ID:    metric.Id,
				MType: "gauge",
				Value: value,
			}
			metrics = append(metrics, m)

		case "counter":
			delta := &metric.Delta
			var m = models.Metrics{
				ID:    metric.Id,
				MType: "counter",
				Delta: delta,
			}
			metrics = append(metrics, m)

		default:
			return errors.New("wrong Metrics type")
		}

	}
}
