package server

import "fmt"

type Storager interface {
	AddGauge(mName string, mValue float64)
	AddCounter(mName string, mValue int64)
	ReadAllGauge() map[string]float64
	ReadAllCounter() map[string]int64
	ReadGauge(mName string) (float64, bool)
	ReadCounter(mName string) (int64, bool)
}

// Storager
// update gauge
func (s *ServerService) UpdateGauge(mName string, mValue float64) error {
	if mName == "" {
		return errMetricNameIncorrect
	}

	s.Storage.AddGauge(mName, mValue)
	if s.cfg.Interval == 0 {
		s.Backup()
	}
	return nil
}

// оupdate counter
func (s *ServerService) UpdateCounter(mName string, mValue int64) error {
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
	if s.cfg.Interval == 0 {
		s.Backup()
	}
	//no error
	return nil
}

// get gauge metric
func (s *ServerService) ViewGaugeByName(mName string) (float64, error) {
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
func (s *ServerService) ViewCounterByName(mName string) (int64, error) {
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
func (s *ServerService) MetricPage() string {
	gauges := s.Storage.ReadAllGauge()
	counters := s.Storage.ReadAllCounter()

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
