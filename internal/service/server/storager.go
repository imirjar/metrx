package server

import (
	"fmt"
	"log"
	"strconv"

	"github.com/imirjar/metrx/internal/models"
)

// get all metrics as html page
func (s ServerService) MetricPage() (string, error) {
	gauges, err := s.MemStorager.ReadAllGauges()
	if err != nil {
		return "", err
	}
	counters, err := s.MemStorager.ReadAllCounters()
	if err != nil {
		return "", err
	}

	gaugeForm := "<a>Gauge</a>"
	for k, g := range gauges {
		gaugeForm += fmt.Sprintf("<li>%s:%f</li>", k, g)
	}

	counterForm := "<a>Counter</a>"
	for k, c := range counters {
		counterForm += fmt.Sprintf("<li>%s:%d</li>", k, c)
	}

	form := fmt.Sprintf("<html><ul>%s</ul><ul>%s</ul></html>", gaugeForm, counterForm)

	return form, nil
}

func (s ServerService) View(metric models.Metrics) (models.Metrics, error) {
	switch metric.MType {
	case "gauge":
		value, ok := s.MemStorager.ReadGauge(metric.ID)
		log.Println("###GAUGE OUT--->", metric.ID, ":", value)
		if !ok {
			return metric, errServiceError
		}
		metric.Value = &value
		return metric, nil
	case "counter":
		delta, ok := s.MemStorager.ReadCounter(metric.ID)
		log.Println("###GAUGE OUT--->", metric.ID, ":", delta)
		if !ok {
			return metric, errServiceError
		}
		metric.Delta = &delta
		return metric, nil
	default:
		return metric, errServiceError
	}
}

func (s ServerService) Update(metric models.Metrics) (models.Metrics, error) {

	switch metric.MType {
	case "gauge":
		if metric.Value == nil {
			return metric, errServiceError
		}
		value, err := s.MemStorager.AddGauge(metric.ID, *metric.Value)
		if err != nil {
			return metric, err
		}
		metric.Value = &value
		return metric, nil

	case "counter":
		if metric.Delta == nil {
			return metric, errServiceError
		}
		delta, err := s.MemStorager.AddCounter(metric.ID, *metric.Delta)
		if err != nil {
			return metric, err
		}
		metric.Delta = &delta
		return metric, nil
	default:
		return metric, errServiceError
	}
}

func (s ServerService) ViewPath(name, mType string) (string, error) {
	switch mType {
	case "gauge":
		value, ok := s.MemStorager.ReadGauge(name)
		log.Println("###GAUGE OUT--->", name, ":", value)
		if !ok {
			return "", errServiceError
		}
		return fmt.Sprint(value), nil
	case "counter":
		delta, ok := s.MemStorager.ReadCounter(name)
		log.Println("###GAUGE OUT--->", name, ":", delta)
		if !ok {
			return "", errServiceError
		}
		return fmt.Sprint(delta), nil
	default:
		return "", errServiceError
	}
}

func (s ServerService) UpdatePath(name, mType, mValue string) (string, error) {

	switch mType {
	case "gauge":
		value, err := strconv.ParseFloat(mValue, 64)
		if err != nil {
			return "", errServiceError
		}
		_, err = s.MemStorager.AddGauge(name, value)
		if err != nil {
			return "", err
		}
		return mValue, nil

	case "counter":
		delta, err := strconv.ParseInt(mValue, 10, 64)
		if err != nil {
			return "", errServiceError
		}
		newDelta, err := s.MemStorager.AddCounter(name, delta)
		if err != nil {
			return "", err
		}

		return fmt.Sprint(newDelta), nil
	default:
		return "", errServiceError
	}
}

func (s ServerService) BatchUpdate(metrics []models.Metrics) error {
	var (
		gauges   = map[string]float64{}
		counters = map[string]int64{}
	)

	for _, metric := range metrics {
		switch metric.MType {
		case "gauge":
			// log.Println("###GAUGE IN--->", metric.ID, ":", *metric.Value)
			gauges[metric.ID] = *metric.Value
		case "counter":
			// log.Println("###COUNTER IN--->", metric.ID, ":", *metric.Delta)
			if _, ok := counters[metric.ID]; ok {
				counters[metric.ID] = counters[metric.ID] + *metric.Delta
			} else {
				counters[metric.ID] = *metric.Delta
			}

		}
	}
	log.Println("###GAUGE IN--->", gauges)
	log.Println("###COUNTER IN--->", counters)
	err := s.MemStorager.AddGauges(gauges)
	if err != nil {
		// log.Println("###GAUGE ERR--->", err)
		return err
	}

	err = s.MemStorager.AddCounters(counters)
	if err != nil {
		// log.Println("###COUNTER ERR--->", err)
		return err
	}

	return nil
}
