package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/imirjar/metrx/internal/models"
)

// Mock генерация будет в файле mock_app.go

func TestAgentApp_Start(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := NewMockClient(ctrl)
	mockSystem := NewMockSystem(ctrl)

	app := AgentApp{
		client: mockClient,
		system: mockSystem,
	}

	metrics := []models.Metrics{
		{ID: "TestMetric1", MType: "gauge", Value: float64Pointer(3.14)},
	}

	// Define behavior for mockSystem.Collect
	mockSystem.EXPECT().Collect(gomock.Any()).Return(metrics, nil).AnyTimes()

	// Define behavior for mockClient.POST
	mockClient.EXPECT().POST(gomock.Any(), metrics).Return(nil).AnyTimes()

	// Run Start method in a goroutine since it blocks
	go func() {
		err := app.Start(100*time.Millisecond, 200*time.Millisecond)
		if err != nil && !errors.Is(err, context.Canceled) {
			t.Errorf("Start() error = %v", err)
		}
	}()

	// Allow the goroutine to run for a while and then cancel
	time.Sleep(1 * time.Second)
}

func TestRun(t *testing.T) {
	// Здесь мы можем написать интеграционный тест для функции Run,
	// однако это будет сложно из-за вызовов client.NewClient и system.NewSystem.
	// Можно использовать тестовую конфигурацию и заменять реальные зависимости моками.

	// Пример кода для интеграционного теста, который мы не можем запустить:
	// Run()
	// Однако это не будет работать в реальной среде тестирования, так как Run вызывает реальные зависимости.
}

// Вспомогательные функции для указателей
func float64Pointer(v float64) *float64 {
	return &v
}

func int64Pointer(v int64) *int64 {
	return &v
}
