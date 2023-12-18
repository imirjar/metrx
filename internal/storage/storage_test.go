package storage

import (
	"testing"

	"github.com/imirjar/metrx/internal/models"
)

func TestAddCounter(t *testing.T) {
	memstorage := New()
	tests := []struct { // добавляем слайс тестов
		name  string
		value models.Counter
		want  int64
	}{
		{
			name: "Create",
			value: models.Counter{
				Name:  "metric",
				Value: 10,
			},
			want: 10,
		},
		{
			name: "Update",
			value: models.Counter{
				Name:  "metric",
				Value: 10,
			},
			want: 20,
		},
		{
			name: "Update",
			value: models.Counter{
				Name:  "metric",
				Value: -20,
			},
			want: 0,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if add, _ := memstorage.AddCounter(test.value); add.Value != test.want {
				t.Errorf("Sum() = %d, want %d", add.Value, test.want)
			}
		})
	}
}
