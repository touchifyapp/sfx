package main

import (
	"log"
)

func main() {
	err := parseArguments()
	if err != nil {
		log.Fatal(err)
	}

	err = bundle()
	if err != nil {
		log.Fatal(err)
	}
}

func bundle() error {
	exefile, exeinfo, err := openExe(args.Exe)
	if err != nil {
		return err
	}
	defer exefile.Close()

	size, err := appendConfigFile(exefile)
	if err != nil {
		return err
	}

	err = appendZipFile(exefile, exeinfo.Size()+size)
	if err != nil {
		return err
	}

	return nil
}
