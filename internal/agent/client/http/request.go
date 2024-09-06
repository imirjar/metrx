package http

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"

	"github.com/imirjar/metrx/internal/models"
	"github.com/imirjar/metrx/pkg/encrypt"
)

func New(metrics []models.Metrics) *Request {
	body, err := json.Marshal(metrics)
	if err != nil {
		return nil
	}

	var request = Request{
		Headers: map[string]string{
			"Content-Type":     "application/json",
			"Content-Encoding": "gzip",
			"X-Real-IP":        "",
		},
	}
	request.Body.Write(body)

	return &request
}

type Request struct {
	Headers map[string]string
	Body    bytes.Buffer
}

func (r *Request) Push(ctx context.Context, url string) error {
	// log.Println(string(r.Body))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url+"/updates", &r.Body)
	if err != nil {
		return err
	}

	for h, i := range r.Headers {
		req.Header.Add(h, i)
	}

	client := http.DefaultClient

	resp, err := client.Do(req)
	if err != nil {
		log.Println("ERROR:", err)
		return err
	}
	defer resp.Body.Close()
	log.Print(resp.Status)

	return nil
}

func (r *Request) Hash(secret string) error {
	var body string
	r.Body.Write([]byte(body))
	hash := encrypt.EncryptSHA256(body, secret)
	r.Headers["HashSHA256"] = hex.EncodeToString(hash)
	// log.Println("##", hex.EncodeToString(hash))
	return nil
}

func (r *Request) Encrypt(pk *rsa.PublicKey) error {
	rnd := rand.Reader
	var body []byte
	r.Body.Write(body)
	encBody, err := rsa.EncryptPKCS1v15(rnd, pk, body)
	if err != nil {
		return err
	}
	r.Body.Read(encBody)

	// log.Println(r.body)
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
