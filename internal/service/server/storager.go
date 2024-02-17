package server

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/imirjar/metrx/internal/models"
)

// get all metrics as html page
func (s ServerService) MetricPage() string {
	gauges := s.MemStorager.ReadAllGauge()
	counters := s.MemStorager.ReadAllCounter()

	gaugeForm := "<a>Gauge</a>"
	for i, g := range gauges {
		gaugeForm += fmt.Sprintf("<li>%s:%f</li>", i, g)
	}

	counterForm := "<a>Counter</a>"
	for i, c := range counters {
		counterForm += fmt.Sprintf("<li>%s:%d</li>", i, c)
	}

	form := fmt.Sprintf("<html><ul>%s</ul><ul>%s</ul></html>", gaugeForm, counterForm)
	if s.cfg.Interval == 0 {
		s.Backup()
	}
	return form
}

func (s ServerService) ByteUpdate(bMetric []byte) ([]byte, error) {
	var metric models.Metrics
	if err := json.Unmarshal(bMetric, &metric); err != nil {
		return nil, err
	}

	switch metric.MType {
	case "gauge":
		s.MemStorager.AddGauge(metric.ID, *metric.Value)
		return bMetric, nil
	case "counter":
		newDelta, err := s.MemStorager.AddCounter(metric.ID, *metric.Delta)
		if err != nil {
			return nil, err
		}
		metric.Delta = &newDelta
		return bMetric, nil
	default:
		return nil, errServiceError
	}
}

func (s ServerService) ByteRead(bMetric []byte) ([]byte, error) {
	var metric models.Metrics
	if err := json.Unmarshal(bMetric, &metric); err != nil {
		return nil, err
	}

	switch metric.MType {
	case "gauge":
		value, ok := s.MemStorager.ReadGauge(metric.ID)
		if !ok {
			return nil, errServiceError
		}
		metric.Value = &value

	case "counter":
		delta, ok := s.MemStorager.ReadCounter(metric.ID)
		if !ok {
			return nil, errServiceError
		}
		metric.Delta = &delta

	default:
		return nil, errServiceError
	}

	b, err := json.Marshal(metric)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return nil, errServiceError
	}

	return b, nil
}

func (s ServerService) Update(mName, mType, mValue string) error {

	switch mType {
	case "gauge":
		value, err := strconv.ParseFloat(mValue, 64)
		if err != nil {
			return errConvertationError
		}
		_, err = s.MemStorager.AddGauge(mName, value)
		return err
	case "counter":
		delta, err := strconv.ParseInt(mValue, 10, 64)
		if err != nil {
			return errConvertationError
		}
		_, err = s.MemStorager.AddCounter(mName, delta)
		return err
	default:
		return errServiceError
	}
}

func (s ServerService) View(mName, mType string) (string, error) {
	switch mType {
	case "gauge":
		v, ok := s.MemStorager.ReadGauge(mName)
		if !ok {
			return "", errServiceError
		}
		value := fmt.Sprintf("%g", v)
		return value, nil
	case "counter":
		d, ok := s.MemStorager.ReadCounter(mName)
		if !ok {
			return "", errServiceError
		}
		delta := fmt.Sprintf("%d", d)
		return delta, nil
	default:
		return "", errServiceError
	}
}
