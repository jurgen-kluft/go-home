package config

import (
	"reflect"
	"testing"
)

func TestPresenceConfigFromJSON(t *testing.T) {
	type args struct {
		jsonstr string
	}
	tests := []struct {
		name    string
		args    args
		want    *PresenceConfig
		wantErr bool
	}{
		{
			name: "simple",
			args: args{jsonstr: "{ \"name\": \"netgear\", \"host\": \"10.0.0.3\", \"port\": 5000, \"user\": \"admin\", \"password\": \"*********\", \"update_history\": 4, \"update_interval_sec\": 10 }"},
			want: &PresenceConfig{
				Name:              "netgear",
				Host:              "10.0.0.3",
				Port:              5000,
				User:              "admin",
				Password:          "*********",
				UpdateHistory:     4,
				UpdateIntervalSec: 10,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PresenceConfigFromJSON([]byte(tt.args.jsonstr))
			if (err != nil) != tt.wantErr {
				t.Errorf("PresenceConfigFromJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PresenceConfigFromJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
