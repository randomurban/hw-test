package main

import (
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	c := exec.Command(cmd[0], cmd[1:]...)

	c.Stdout = os.Stdout
	c.Stdin = os.Stdin
	c.Stderr = os.Stderr

	for s, value := range env {
		if !value.NeedRemove {
			sv := fmt.Sprintf("%s=%s", s, value.Value)
			c.Env = append(c.Env, sv)
		}
	}
	err := c.Run()
	if err != nil {
		fmt.Println(err.Error())
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode()
		} else {
			return 111
		}
	}
	return 0
}
