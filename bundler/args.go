package main

import (
	"flag"
)

var args struct {
	ID      string
	Run     string
	Dest    string
	Args    string
	Version string

	Dir     string
	Exe     string
	Verbose bool
}

func parseArguments() error {
	id := flag.String("id", "co.touchify.sfx", "The unique ID for this package.")
	run := flag.String("run", "", "The program to run in the project directory (default: auto-detect).")
	dest := flag.String("dest", "", "The absolute destination path to extract project in (default: temp).")
	cargs := flag.String("args", "", "arguments to pass to executable")
	version := flag.String("version", "1.0.0", "The program to run in the project directory.")

	dir := flag.String("dir", "project", "The directory to bundle into sfx.")
	exe := flag.String("exe", "sfx.exe", "The program to bundle the project in.")
	verb := flag.Bool("v", false, "Enable program output.")

	flag.Parse()

	args.ID = *id
	args.Run = *run
	args.Dest = *dest
	args.Args = *cargs
	args.Version = *version

	args.Dir = *dir
	args.Exe = *exe
	args.Verbose = *verb

	if args.Run == "" {
		run, err := findExeInDir(args.Dir)
		if err != nil {
			return err
		}

		args.Run = run
	}

	return nil
}
