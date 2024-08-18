package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go-envdir /path/to/env/dir command arg1 arg2 ...")
		os.Exit(111)
	}
	env, err := ReadDir(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(111)
	}
	returnCode := RunCmd(os.Args[2:], env)
	os.Exit(returnCode)
}
