package client

import (
	"context"
	"crypto/rsa"
	"log"
	"sync"

	"github.com/imirjar/metrx/internal/agent/client/grpc"
	"github.com/imirjar/metrx/internal/models"
	"github.com/imirjar/metrx/pkg/ping"
)

// Client application part witch provide
type Client struct {
	secret string         // secret is using for hash function which place in HashSHA256 HTTP header
	target string         // server ip where we will sent request
	host   string         // our ip for X-Real-IP HTTP header
	pk     *rsa.PublicKey // public key for security
	sync.Mutex
}

func New(secret, target string, crypto *rsa.PublicKey) *Client {

	// get host IP
	host, err := ping.GetMyIP()
	if err != nil {
		log.Print(err)
	}

	// make client
	cli := Client{
		secret: secret,
		host:   host,
		target: target,
		pk:     crypto,
	}

	return &cli
}

// Using http.Client to sent our metrics to host+/updates
func (c *Client) POST(ctx context.Context, batch []models.Metrics) error {

	// var request Requester

	// request := http.New(batch)
	request := grpc.New(batch)

	//Add middlewares
	request.Hash(c.secret)
	request.Encrypt(c.pk)

	return request.Push(ctx, c.target)

}

type Requester interface {
	Push(ctx context.Context, url string) error
	Hash(secret string) error
	Encrypt(pk *rsa.PublicKey) error
}
