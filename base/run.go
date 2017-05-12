package main

import (
	"os"
	"path/filepath"
	"strings"
)

func run(config *peConfig) error {
	exe := filepath.Join(config.Dest, config.Run)

	attr := os.ProcAttr{
		Dir: filepath.Dir(exe),
		Env: os.Environ(),
		Files: []*os.File{
			nil,
			nil,
			nil,
		},
	}

	args := []string{
		exe,
	}

	if config.Args != "" {
		iargs := strings.Split(config.Args, " ")
		args = append(args, iargs...)
	}

	procArgsLen := len(os.Args)
	if procArgsLen > 1 {
		for i := 1; i < procArgsLen; i++ {
			args = append(args, os.Args[i])
		}
	}

	process, err := os.StartProcess(exe, args, &attr)
	if err != nil {
		return err
	}

	err = process.Release()
	return err
}
