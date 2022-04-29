package main

import (
	"os"
)

func Process(cmd []string) {
	if len(cmd) < 3 {
		os.Exit(-1)
	}

	envmap, err := ReadDir(cmd[1])
	if err != nil {
		println("err: ", err)
		os.Exit(-1)
	}

	rc := RunCmd(cmd[2:], envmap)
	os.Exit(rc)
}
