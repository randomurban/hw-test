package main

import "testing"

func TestRunCmd(t *testing.T) {
	type args struct {
		cmd []string
		env Environment
	}
	tests := []struct {
		name           string
		args           args
		wantReturnCode int
	}{
		{
			"Simple command",
			args{
				[]string{"/bin/bash", "-c", "./testdata/runme.sh", "arg1=1", "arg2=2"},
				Environment{
					"HELLO": {"bar", false},
					"FOO":   {"   foo\nwith new line", false},
					"UNSET": {"", true},
					"ADDED": {"from original env", false},
					"EMPTY": {"", true},
				},
			},
			0,
		},

		//export HELLO="SHOULD_REPLACE"
		//export FOO="SHOULD_REPLACE"
		//export UNSET="SHOULD_REMOVE"
		//export ADDED="from original env"
		//export EMPTY="SHOULD_BE_EMPTY"
		//
		//result=$(./go-envdir "$(pwd)/testdata/env" "/bin/bash" "$(pwd)/testdata/echo.sh" arg1=1 arg2=2)
		//expected='HELLO is ("hello")
		//BAR is (bar)
		//FOO is (   foo
		//with new line)
		//UNSET is ()
		//ADDED is (from original env)
		//EMPTY is ()
		//arguments are arg1=1 arg2=2'

		// windows
		//{"Simple command", args{[]string{"ls"}, Environment{"TESTDIR": {"testdata/env", false}}}, 111},
		//{"Windows Simple command", args{[]string{"cmd.exe", "/c", "dir", "%TESTDIR%"}, Environment{"TESTDIR": {"testdata\\env", false}}}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotReturnCode := RunCmd(tt.args.cmd, tt.args.env); gotReturnCode != tt.wantReturnCode {
				t.Errorf("RunCmd() = %v, want %v", gotReturnCode, tt.wantReturnCode)
			}
		})
	}
}
