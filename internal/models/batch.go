package models

type Batch struct {
	Metrics []Metrics
}

func (b *Batch) AddGauge(name string, value float64) {
	metric := Metrics{
		ID:    name,
		MType: "gauge",
		Value: &value,
	}
	// log.Println("#####batch.go value-->", value)
	b.Metrics = append(b.Metrics, metric)
}

func (b *Batch) AddCounter(name string, delta int64) {
	metric := Metrics{
		ID:    name,
		MType: "counter",
		Delta: &delta,
	}
	// log.Println("#####batch.go delta-->", delta)
	b.Metrics = append(b.Metrics, metric)
}
