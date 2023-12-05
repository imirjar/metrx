package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strconv"
)

type Gauge struct {
	Name  string  `json:name`
	Value float64 `json:value`
}

type Counter struct {
	Name  string  `json:name`
	Value float64 `json:value`
}

type MemStorage struct {
	CounterStorage []Counter
	GaugeStorage   []Gauge
}

var store MemStorage

func gaugeHandle(w http.ResponseWriter, r *http.Request) {

	dir, file := path.Split(r.URL.Path)
	name := path.Base(dir)

	value, err := strconv.ParseFloat(file, 64)
	if err != nil {
		panic(err)
	}

	gauge := Gauge{
		Name:  name,
		Value: value,
	}

	store.GaugeStorage = append(store.GaugeStorage, gauge)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(store.GaugeStorage)
}

func counterHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("counter"))
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	gauge := http.NewServeMux()
	gauge.Handle("/", middleware(http.HandlerFunc(gaugeHandle)))

	counter := http.NewServeMux()
	counter.Handle("/", middleware(http.HandlerFunc(counterHandle)))

	update := http.NewServeMux()
	update.Handle("/gauge/", http.StripPrefix("/gauge", gauge))
	update.Handle("/counter/", http.StripPrefix("/counter", counter))

	mux := http.NewServeMux()
	mux.Handle("/update/", http.StripPrefix("/update", update))

	return http.ListenAndServe(":8080", mux)
}

func middleware(next http.Handler) http.Handler {

	// получаем Handler приведением типа http.HandlerFunc
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// здесь пишем логику обработки
		// например, разрешаем запросы cross-domain
		// w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "text/plain")
		// замыкание: используем ServeHTTP следующего хендлера
		next.ServeHTTP(w, r)
	})
}
