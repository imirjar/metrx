package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSetval(t *testing.T) {

	type sent struct {
		metric Metrics
		value  string
	}

	type want struct {
		err   error
		value string
	}

	tests := []struct {
		name string
		want want
		sent sent
	}{
		{
			name: "ok gauge",
			sent: sent{
				metric: Metrics{
					ID:    "Gauge",
					MType: "gauge",
				},
				value: "10.999999999999998",
			},
			want: want{
				err:   nil,
				value: "10.999999999999998",
			},
		},
		{
			name: "ok counter",
			sent: sent{
				metric: Metrics{
					ID:    "Counter",
					MType: "counter",
				},
				value: "100",
			},
			want: want{
				err:   nil,
				value: "100",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sent.metric.SetVal(tt.sent.value)
			if err != nil {
				t.Error(err)
			}

			value, err := tt.sent.metric.GetVal()
			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, tt.want.err, err)
			assert.Equal(t, tt.sent.value, value)

		})
	}

}
