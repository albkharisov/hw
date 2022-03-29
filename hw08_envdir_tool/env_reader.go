package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	MaxEnvVarSize = 512
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
// Skip files which are > MaxEnvVarSize (512).
func ReadDir(dir string) (Environment, error) {
	envmap := make(Environment)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, errWalk error) error {
		if errWalk != nil || !info.Mode().IsRegular() {
			if errWalk != nil {
				fmt.Println(errWalk)
				return errWalk
			}
			return nil
		}

		if info.Size() > MaxEnvVarSize {
			return nil
		}

		file, errOpen := os.Open(path)
		if errOpen != nil {
			return nil
		}
		defer file.Close()

		var value EnvValue
		if info.Size() == 0 {
			value.NeedRemove = true
		} else {
			reader := bufio.NewReader(file)
			line, _, err := reader.ReadLine()
			if err != nil {
				return nil
			}

			line = bytes.ReplaceAll(line, []byte{0x00}, []byte{'\n'})
			value.NeedRemove = false
			value.Value = strings.TrimRight(string(line), " \t")
		}
		envmap[filepath.Base(path)] = value

		return nil
	})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return envmap, nil
}
