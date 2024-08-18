package main

import (
	"reflect"
	"testing"
)

func TestReadDir(t *testing.T) {
}

func Test_newEnv(t *testing.T) {
	type args struct {
		osEnv []string
	}
	tests := []struct {
		name string
		args args
		want Environment
	}{
		{
			"Empty osEnv",
			args{[]string{}},
			Environment{},
		},
		{
			"Simple osEnv",
			args{[]string{"TESTENV=test_env"}},
			Environment{"TESTENV": {"test_env", false}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newEnv(tt.args.osEnv); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}
