package service

import (
	"fmt"
)

type Server struct {
	Storage Storager
}

// update gauge
func (s *Server) UpdateGauge(mName string, mValue float64) error {
	if mName == "" {
		return errMetricNameIncorrect
	}

	s.Storage.AddGauge(mName, mValue)
	return nil
}

// оupdate counter
func (s *Server) UpdateCounter(mName string, mValue int64) error {
	if mName == "" {
		return errMetricNameIncorrect
	}

	// if counter exists -> counter += new value
	curV, ok := s.Storage.ReadCounter(mName)
	if ok {
		s.Storage.AddCounter(mName, curV+mValue)
	} else {
		s.Storage.AddCounter(mName, mValue)
	}

	//no error
	return nil
}

// get gauge metric
func (s *Server) ViewGaugeByName(mName string) (float64, error) {
	if mName == "" {
		return 0, errMetricNameIncorrect
	}

	gauge, ok := s.Storage.ReadGauge(mName)
	if !ok {
		return gauge, errServiceError
	}

	return gauge, nil
}

// get counter metric
func (s *Server) ViewCounterByName(mName string) (int64, error) {
	if mName == "" {
		return 0, errMetricNameIncorrect
	}

	counter, ok := s.Storage.ReadCounter(mName)
	if !ok {
		return counter, errServiceError
	}

	return counter, nil

}

// get all metrics as html page
func (s *Server) MetricList() string {
	metric := s.Storage.ReadAll()

	gaugeForm := "<a>Gauge</a>"
	for i, g := range metric.Gauge {
		gaugeForm += fmt.Sprintf("<li>%s:%f</li>", i, g)
	}

	counterForm := "<a>Counter</a>"
	for i, c := range metric.Counter {
		counterForm += fmt.Sprintf("<li>%s:%d</li>", i, c)
	}

	form := fmt.Sprintf("<html><ul>%s</ul><ul>%s</ul></html>", gaugeForm, counterForm)
	return form
}
