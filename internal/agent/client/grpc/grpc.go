package grpc

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	"io"
	"log"

	"github.com/imirjar/metrx/internal/api"
	"github.com/imirjar/metrx/internal/models"
	"github.com/imirjar/metrx/pkg/encrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func New(batch []models.Metrics) *Request {
	return &Request{
		Metrcis: batch,
		ConOps:  []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
		Headers: map[string]string{},
	}
}

// Request application part witch provide
type Request struct {
	Headers map[string]string
	ConOps  []grpc.DialOption
	Metrcis []models.Metrics
}

// Function make connection create request and send stream
func (r *Request) Push(ctx context.Context, url string) error {

	// Create Connection
	conn, err := grpc.NewClient(url, r.ConOps...) //fmt.Sprintf("localhost:%s", "3200")
	if err != nil {
		log.Println("ERROR because the GRPC Client was not created", err)
		return err
	}
	defer conn.Close()

	// Create proto client
	client := api.NewGoMetricsClient(conn)

	// Create grpc metrics list
	var metrics []*api.Metric
	for _, m := range r.Metrcis {

		switch m.MType {
		case "gauge":
			metric := api.Metric{
				Id:    m.ID,
				Type:  "gauge",
				Value: *m.Value,
			}
			// log.Print(metric.Value)
			metrics = append(metrics, &metric)
		case "counter":
			metric := api.Metric{
				Id:    m.ID,
				Type:  "counter",
				Delta: *m.Delta,
			}
			metrics = append(metrics, &metric)
		default:
			log.Println("WTF METRIC!??!?!?!?!?!?")
		}
	}

	// Create stream for metric list
	stream, err := client.BatchUpdate(ctx)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
		return err
	}

	for _, metric := range metrics {
		if err := stream.Send(metric); err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("%v.Send(%v) = %v", stream, metric, err)
		}
	}

	//Adding headers
	for ht, hb := range r.Headers {
		log.Printf("\n =========== \n Заголовок: %s \n Значение: %s \n ===========", ht, hb)
	}
	// md := metadata.Pairs("X-Real-IP", c.host)

	// Catch request
	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
	}

	log.Printf("Route summary: %v", reply)
	return nil
}

// Adding hash sum to header
func (r *Request) Hash(secret string) error {
	hash := encrypt.EncryptSHA256(fmt.Sprint(&r.Metrcis), secret)
	r.Headers["HashSHA256"] = hex.EncodeToString(hash)
	return nil
}

func (r *Request) Encrypt(pk *rsa.PublicKey) error {
	rnd := rand.Reader
	var body []byte

	encBody, err := rsa.EncryptPKCS1v15(rnd, pk, body)
	if err != nil {
		return err
	}

	log.Print(encBody)

	// opt := grpc.WithTransportCredentials(encBody)
	// r.ConOps = append(r.ConOps, opt)

	return nil
}

func (r *Request) Gzip() error {
	// encBody, err := rsa.EncryptPKCS1v15(rnd, pk, r.Body)
	// if err != nil {
	// 	return err
	// }
	// r.Body = cmpBody
	return nil
}
