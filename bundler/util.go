package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func openExe(exe string) (*os.File, os.FileInfo, error) {
	exeAbs, err := filepath.Abs(exe)
	if err != nil {
		return nil, nil, err
	}

	exefile, err := os.OpenFile(exeAbs, os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return nil, nil, err
	}

	exefileInfo, err := exefile.Stat()
	if err != nil {
		return nil, nil, err
	}

	return exefile, exefileInfo, nil
}

func findExeInDir(dir string) (string, error) {
	exe := ""
	dirAbs, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}

	verbosef("No run flag passed, trying to find executable in %s", dirAbs)

	err = filepath.Walk(dirAbs, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && path != dirAbs {
			return filepath.SkipDir
		}

		if exe == "" && filepath.Ext(path) == ".exe" {
			exe = strings.TrimPrefix(path, dirAbs+string(filepath.Separator))
			verbosef("Found executable: %s (%s)", path, exe)
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	if exe == "" {
		return "", errors.New("No executable found in dir")
	}

	return exe, nil
}

func verbosef(format string, stuff ...interface{}) {
	if args.Verbose {
		log.Printf(format, stuff...)
	}
}
