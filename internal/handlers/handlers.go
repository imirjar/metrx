package handlers

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/imirjar/metrx/internal/service"
)

func parsePath(path string) ([]string, error) {
	params := strings.Split(path, "/")
	if len(params) == 5 {
		for i := range params {
			if params[i] == "" {
				return nil, fmt.Errorf("Incorrect path")
			}
		}
		return params, nil
	} else {
		return nil, fmt.Errorf("Incorrect path")
	}
}

type Handler struct {
	Service service.Service
}

func New() *Handler {

	return &Handler{
		// Routes:  defineRoutes(),
		Service: *service.New(),
	}
}

func (h *Handler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	// params, err := parsePath(r.URL.Path)
	// if err != nil {
	// 	fmt.Println(err)
	// 	w.WriteHeader(http.StatusNotFound)
	// 	w.Write([]byte("Неверный запрос"))
	// 	return
	// }

	// switch metrix := params[1]; metrix {
	// case "gauge":
	// 	fmt.Println("gauge")
	// 	gauge := h.Service.Gauge(params[2], params[3])
	// 	w.WriteHeader(http.StatusOK)
	// 	json.NewEncoder(w).Encode(gauge)
	// case "counter":
	// 	fmt.Println("counter")
	// 	counter := h.Service.Counter(params[2], params[3])

	// 	w.WriteHeader(http.StatusOK)
	// 	json.NewEncoder(w).Encode(counter)
	// default:
	// 	w.WriteHeader(http.StatusNotFound)
	// 	w.Write([]byte("Неверный запрос"))
	// 	return
	// }
	// metric := path.Dir(r.URL.Path)
	// fmt.Println(metric)

	base, value := path.Split(r.URL.Path)
	fmt.Println(value, "###", base)

	secondBase, metrixName := path.Split(base)
	fmt.Println(secondBase, "###", metrixName)
	// name := path.Base(path.Dir(r.URL.Path))
	// fmt.Println(name)

	w.Write([]byte(value))
}

// func (h *Handler) GaugeHandle(w http.ResponseWriter, r *http.Request) {
// 	params := strings.Split(r.URL.Path, "/")
// 	fmt.Println(params)
// 	gauge := h.Service.Gauge(params[1], params[2])

// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(gauge)
// }

// func (h *Handler) CounterHandle(w http.ResponseWriter, r *http.Request) {

// 	namePath, value := path.Split(r.URL.Path)
// 	name := path.Base(namePath)

// 	counter := h.Service.Counter(name, value)

// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(counter)
// }

func (h *Handler) DefineRoutes() *http.ServeMux {

	mux := http.NewServeMux()
	mux.Handle("/update/", http.HandlerFunc(h.UpdateHandler))

	return mux
}
