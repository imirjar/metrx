package grpc

import (
	"context"
	"log"

	"github.com/imirjar/metrx/internal/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/internal/metadata"
)

// Client application part witch provide
type Client struct {
	addr string
	host string
}

func New(host, addr string) *Client {
	return &Client{
		addr: addr,
		host: host,
	}
}

func (c *Client) POST(ctx context.Context, body []byte) error {
	conn, err := grpc.NewClient(c.addr) //fmt.Sprintf("localhost:%s", "3200")
	if err != nil {
		log.Print(err)
	}
	defer conn.Close()

	cli := api.NewGoMetricsClient(conn)

	md := metadata.Pairs("X-Real-IP", c.host)

	r, err := cli.BatchUpdate(ctx, &api.Request{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetResponse())
	return nil
}
