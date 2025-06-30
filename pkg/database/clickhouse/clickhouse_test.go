package clickhouse

import (
	"context"
	"testing"
)

func TestConnectClickhouse(t *testing.T) {
	type args struct {
		ctx  context.Context
		cfg  *ClickhouseConfig
		want *DatabaseClickhouse
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "",
			args: args{
				ctx: context.Background(),
				cfg: &ClickhouseConfig{
					Host:     "127.0.0.1",
					Port:     9000,
					Username: "test",
					Password: "test123",
					DBName:   "gds",
				},
				want: &DatabaseClickhouse{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConnectClickhouse(tt.args.ctx, tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConnectClickhouse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got.db == nil {
				t.Errorf("ConnectClickhouse() got = %v, want %v", got, tt.args.want)
			}
		})
	}
}
