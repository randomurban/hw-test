package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	res := newEnv(os.Environ())
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			name := entry.Name()
			value, err := read(filepath.Join(dir, name))
			if err != nil {
				return nil, err
			}
			res[name] = value
		}
	}
	return res, nil
}

func read(name string) (EnvValue, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return EnvValue{}, err
	}
	if len(data) == 0 {
		return EnvValue{"", true}, nil
	}
	dataLine, _, _ := strings.Cut(string(data), "\n")
	dataZero := bytes.Replace([]byte(dataLine), []byte{0}, []byte("\n"), -1)
	s := strings.TrimRight(string(dataZero), " \t\r")
	if len(s) == 0 {
		return EnvValue{"", true}, nil
	}
	return EnvValue{Value: s, NeedRemove: false}, nil
}

func newEnv(osEnv []string) Environment {
	res := make(Environment)
	for _, s := range osEnv {
		key, val, _ := strings.Cut(s, "=")
		res[key] = EnvValue{Value: val, NeedRemove: false}
	}
	return res
}
