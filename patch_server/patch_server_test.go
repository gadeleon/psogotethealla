package patchserver

import (
	"net"
	"reflect"
	"testing"
)

func Test_parseIPString(t *testing.T) {
	expected_ip_parse := net.ParseIP("127.0.0.1")
	type args struct {
		ip string
	}
	tests := []struct {
		name    string
		args    args
		want    net.IP
		wantErr bool
	}{
		// TODO: Add test cases.
		{"localhost parsed", args{"localhost"}, expected_ip_parse, false},
		{"127.0.0.1 parsed", args{"127.0.0.1"}, expected_ip_parse, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseIPString(tt.args.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseIPString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseIPString() = %v, want %v", got, tt.want)
			}
		})
	}
}
