package limiter

import (
	"reflect"
	"testing"
	"time"
)

func TestNewGlobalRate(t *testing.T) {
	type args struct {
		command string
		limit   int
	}
	tests := []struct {
		name    string
		args    args
		want    GlobalRate
		wantErr bool
	}{
		{
			name: "normal, with 24 * hour",
			args: args{
				command: "24-h",
				limit:   200,
			},
			want: GlobalRate{
				Command: "24-h",
				Period:  24 * time.Hour,
				Limit:   200,
			},
			wantErr: false,
		},
		{
			name: "normal, with H",
			args: args{
				command: "24-H",
				limit:   200,
			},
			want: GlobalRate{
				Command: "24-H",
				Period:  24 * time.Hour,
				Limit:   200,
			},
			wantErr: false,
		},
		{
			name: "normal, with 20 * minute",
			args: args{
				command: "20-m",
				limit:   250,
			},
			want: GlobalRate{
				Command: "20-m",
				Period:  20 * time.Minute,
				Limit:   250,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGlobalRate(tt.args.command, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGlobalRate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGlobalRate() = %v, want %v", got, tt.want)
			}
		})
	}
}
