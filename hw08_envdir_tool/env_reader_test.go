package main

import (
	"reflect"
	"testing"
)

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

func Test_read(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    EnvValue
		wantErr bool
	}{
		{"Error on Empty name", args{""}, EnvValue{}, true},
		{"Error on wrong file name", args{"some_strange_file"}, EnvValue{}, true},
		{"Multiline File", args{"testdata/env/BAR"}, EnvValue{"bar", false}, false},
		{"Empty File", args{"testdata/env/EMPTY"}, EnvValue{"", true}, false},
		{"Zero Coded File", args{"testdata/env/FOO"}, EnvValue{"   foo\nwith new line", false}, false},
		{"Quotes File", args{"testdata/env/HELLO"}, EnvValue{"\"hello\"", false}, false},
		{"Unset multiline File", args{"testdata/env/UNSET"}, EnvValue{"", true}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := read(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("read() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadDir(t *testing.T) {
	type args struct {
		dir  string
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    EnvValue
		wantErr bool
	}{
		{"Empty dir", args{"", ""}, EnvValue{}, true},
		{"BAR in testdata/env dir", args{"testdata/env", "BAR"}, EnvValue{"bar", false}, false},
		{"EMPTY in testdata/env dir", args{"testdata/env", "EMPTY"}, EnvValue{"", true}, false},
		{"FOO in testdata/env dir", args{"testdata/env", "FOO"}, EnvValue{"   foo\nwith new line", false}, false},
		{"HELLO in testdata/env dir", args{"testdata/env", "HELLO"}, EnvValue{"\"hello\"", false}, false},
		{"UNSET in testdata/env dir", args{"testdata/env", "UNSET"}, EnvValue{"", true}, false},
		{"Unknown file in testdata/env dir", args{"testdata/env", "FOOBAR"}, EnvValue{"", false}, false},
		{"Error on name with '='", args{"testdata/env_err", "ENV=ERR"}, EnvValue{"", false}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAll, err := ReadDir(tt.args.dir)
			got := gotAll[tt.args.name]
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.args.name != "" && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadDir() got = %v, want %v", got, tt.want)
			}
		})
	}
}
