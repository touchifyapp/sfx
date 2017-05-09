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
	exefile, err := openExe(args.Exe)
	if err != nil {
		return err
	}
	defer exefile.Close()

	err = appendConfigFile(exefile)
	if err != nil {
		return err
	}

	err = appendZipFile(exefile)
	if err != nil {
		return err
	}

	return nil
}
