package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall" //nolint:typecheck
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmdLine []string, env Environment) (returnCode int) {
	for i := range env {
		if env[i].NeedRemove {
			os.Unsetenv(i)
		} else {
			os.Setenv(i, env[i].Value)
		}
	}

	// #nosec G204
	cmd := exec.Command(cmdLine[0], (cmdLine[1:])...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		var cmdError *exec.ExitError
		if errors.As(err, &cmdError) {
			if status, ok := cmdError.Sys().(syscall.WaitStatus); ok {
				return status.ExitStatus()
			}
		}
		fmt.Printf("err of cmd.Run(): %v\n", err)
		return -1
	}

	return 0
}
