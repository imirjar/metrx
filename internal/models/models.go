package models

type Gauge struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

type Counter struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

type Metrix struct {
	Type  string
	Name  string
	Value any
}
