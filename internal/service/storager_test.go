package service

import (
	"context"
	"testing"

	"github.com/imirjar/metrx/config"
)

func TestServerService_PathHandler(t *testing.T) {

	service := NewServerService(config.Testcfg)
	type args struct {
		mName  string
		mType  string
		mValue string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "OK Gauge",
			args: args{
				mName:  "gaugeMetric",
				mType:  "gauge",
				mValue: "100",
			},
			want:    "100",
			wantErr: false,
		},
		{
			name: "OK Counter",
			args: args{
				mName:  "counterMetric",
				mType:  "counter",
				mValue: "100",
			},
			want:    "100",
			wantErr: false,
		},
		{
			name: "OK Add Counter",
			args: args{
				mName:  "counterMetric",
				mType:  "counter",
				mValue: "100",
			},
			want:    "200",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.UpdatePath(context.Background(), tt.args.mName, tt.args.mType, tt.args.mValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServerService.ViewPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.want {
				t.Errorf("ServerService.ViewPath() = %v, want %v", result, tt.want)
			}
		})
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.ViewPath(context.Background(), tt.args.mName, tt.args.mType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServerService.ViewPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.want {
				t.Errorf("ServerService.ViewPath() = %v, want %v", result, tt.want)
			}
		})
	}
}
