package models

// Batch model use in agent app for collecting metrics from system
type Batch struct {
	Metrics []Metrics
}

// Method trnsform system value in float64 to metric list
func (b *Batch) AddGauge(name string, value float64) {
	metric := Metrics{
		ID:    name,
		MType: "gauge",
		Value: &value,
	}
	b.Metrics = append(b.Metrics, metric)
}

// Method trnsform system value in int64 to metric list
func (b *Batch) AddCounter(name string, delta int64) {
	metric := Metrics{
		ID:    name,
		MType: "counter",
		Delta: &delta,
	}
	b.Metrics = append(b.Metrics, metric)
}
