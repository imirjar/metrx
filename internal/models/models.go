package models

type Gauge struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

type Counter struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

func (c Counter) Sum(i int64) int64 {
	return c.Value + i
}
