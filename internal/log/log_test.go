package log

import (
	"testing"

	"github.com/go-logr/logr"
)

func TestSetLogger(t *testing.T) {
	type args struct {
		log logr.Logger
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetLogger(tt.args.log)
		})
	}
}

func Test_getMsg(t *testing.T) {
	type args struct {
		format string
		args   []interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getMsg(tt.args.format, tt.args.args...); got != tt.want {
				t.Errorf("getMsg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInfof(t *testing.T) {
	type args struct {
		format string
		args   []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "test infof",
			args: args{
				format: "test:%d",
				args:   []interface{}{123},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Infof(tt.args.format, tt.args.args...)
		})
	}
}

func TestErrorf(t *testing.T) {
	type args struct {
		format string
		args   []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "test errorf",
			args: args{
				format: "test:%d",
				args:   []interface{}{123},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Errorf(tt.args.format, tt.args.args...)
		})
	}
}
