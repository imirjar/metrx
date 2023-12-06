package storage

import (
	"fmt"
	"reflect"

	"github.com/imirjar/metrx/internal/models"
)

type MemStorage struct {
	CounterStorage []models.Counter
	GaugeStorage   []models.Gauge
}

var store MemStorage

func New() *MemStorage {
	return &MemStorage{
		CounterStorage: []models.Counter{},
		GaugeStorage:   []models.Gauge{},
	}
}

func (m *MemStorage) Create(obj *models.Gauge) {
	m.GaugeStorage = append(m.GaugeStorage, *obj)
}

func (m *MemStorage) Read(value, field string) *models.Gauge {

	for _, v := range m.GaugeStorage {
		if reflect.ValueOf(v).FieldByName(field).String() == value {
			fmt.Println("Элемент в списке")
			return &v
		}
	}
	fmt.Println("Элемента нет списке")
	return nil
}

func (m *MemStorage) Update(obj *models.Gauge, value float64) {
	fmt.Println("Обновляю")
	obj.Value = value
}
func (m *MemStorage) Delete() {}
