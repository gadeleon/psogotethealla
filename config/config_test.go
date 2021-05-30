package config

import (
	"reflect"
	"testing"

	"github.com/go-ini/ini"
)

func TestConfig_New(t *testing.T) {
	iFile, _ := ini.Load("example.ini")
	validConfig := Config{
		Config: iFile,
	}
	type fields struct {
		Config *ini.File
	}
	type args struct {
		fname string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Config
		wantErr bool
	}{
		// TODO: Add test cases.
		{"Failed to Load", fields{}, args{"no.ini"}, nil, true},
		{"Load Example", fields{}, args{"example.ini"}, &validConfig, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Config: tt.fields.Config,
			}
			got, err := c.New(tt.args.fname)
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.New() = %v, want %v", got, tt.want)
			}
		})
	}
}
