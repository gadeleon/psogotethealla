package patchserver

import (
	"net"
	"reflect"
	"testing"

	"github.com/gadeleon/psogotethealla/config"
	"github.com/go-ini/ini"
)

func Test_convertIPString(t *testing.T) {
	iFile, _ := ini.Load("../config/example.ini")
	validConfig := config.Config{
		Config: iFile,
	}
	type args struct {
		c *config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    net.IP
		wantErr bool
	}{
		// TODO: Add test cases.
		{"localhost parsed", args{&validConfig}, net.ParseIP("127.0.0.1"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertIPString(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertIPString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertIPString() = %v, want %v", got, tt.want)
			}
		})
	}
}
