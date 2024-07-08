package http

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/imirjar/metrx/internal/server/controller/http/middleware/compressor"
	"github.com/imirjar/metrx/internal/server/controller/http/middleware/encryptor"
	"github.com/imirjar/metrx/internal/server/controller/http/middleware/logger"
)

// Http gateway using secret value for encoding
// it has few endpoints for handlers whitch use Service interface
type HTTPGateway struct {
	Service Service
	Server  *http.Server
	pk      *rsa.PrivateKey
}

func NewGateway(path, crypto, secret, conn string) *HTTPGateway {
	gtw := HTTPGateway{}

	if crypto != "" {
		b, err := os.ReadFile(crypto)
		if err != nil {
			log.Print(errLoadPrivateKey)
		}

		block, _ := pem.Decode(b)
		if block == nil || block.Type != "RSA PRIVATE KEY" {
			log.Print("Block type is nil or not RSA PRIVATE KEY", block.Type)
		}

		privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			// log.Print("#####er")
			log.Print(err)
		}

		// rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
		// if !ok {
		// 	log.Print("Ключ не является RSA приватным ключом")
		// }

		gtw.pk = privateKey

	}

	router := chi.NewRouter()
	router.Use(middleware.NoCache)

	router.Use(encryptor.DecryptR(gtw.pk))
	router.Use(compressor.Compressing())
	router.Use(logger.Logger())

	// Save metric
	router.Route("/update", func(update chi.Router) {
		update.Post("/{type}/{name}/{value}", gtw.UpdatePathHandler())
		update.Post("/", gtw.UpdateJSONHandler())
	})

	// Read metric if exists
	router.Route("/value", func(value chi.Router) {
		value.Get("/{type}/{name}", gtw.ValuePathHandler())
		value.Post("/", gtw.ValueJSONHandler())
	})

	// Save list passed metrics
	router.Route("/updates", func(batch chi.Router) {
		batch.Post("/", gtw.BatchHandler())
	})

	// Check db connection
	router.Get("/ping", gtw.Ping(conn))

	// HTML page witch all of metrics
	router.Get("/", gtw.MainPage())

	// Pprof package routes
	router.Mount("/debug", middleware.Profiler())

	gtw.Server = &http.Server{
		Addr:    path,
		Handler: router,
	}

	return &gtw
}
