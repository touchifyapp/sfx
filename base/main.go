package main

import (
	"os"
)

func main() {
	path, err := os.Executable()
	if err != nil {
		verboseFatal(err)
	}

	verbosef("Executable path: %s", path)

	closer, config, reader, err := openExeContent(path)
	if err != nil {
		verboseFatal(err)
	}

	defer closer.Close()

	verbosef("Uncompressing resources to: %s", config.Dest)
	err = uncompress(reader, config)
	if err != nil {
		verboseFatal(err)
	}

	verbosef("Running: %s", config.Run)
	err = run(config)
	if err != nil {
		verboseFatal(err)
	}
}
