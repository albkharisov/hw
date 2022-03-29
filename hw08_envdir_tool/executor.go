package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

//
//func setupEnv(env []string, envMap Environment) []string {
//
//	ret := make([]string, len(env))
//	copy(ret, env)
//
//	for i := range envMap {
//		if envMap[i].NeedRemove {
//			os.Unsetenv(i)
//		} else {
//			ret = append(ret, fmt.Sprintf("%s=%s", i, envMap[i].Value))
//			//			os.Setenv(i, envMap[i].Value)
//		}
//	}
//
//	return ret
//}
//

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmdLine []string, env Environment) (returnCode int) {
	//	os.Clearenv()

	for i := range env {
		if env[i].NeedRemove {
			os.Unsetenv(i)
		} else {
			os.Setenv(i, env[i].Value)
		}
	}

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
