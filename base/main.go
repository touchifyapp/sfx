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

	destConfig, err := readDestConfig(config)
	if err != nil {
		verboseFatal(err)
	}

	mode := getInstallMode(config, destConfig)
	if mode == modOUTDATED {
		verbosef("SFX version (%s) is lower than installed version (%s). Running installed configuration (%s)...", config.Version, destConfig.Version, destConfig.Run)
		err = run(destConfig)
		if err != nil {
			verboseFatal(err)
		}

		return
	}

	switch mode {
	case modEQUAL:
		verbosef("SFX version (%s) is equal to installed version (%s). Checking file dates...", config.Version, destConfig.Version)

	case modUPDATE:
		verbosef("SFX version (%s) is greater than installed version (%s). Force decompression...", config.Version, destConfig.Version)
	}

	verbosef("Uncompressing resources to: %s", config.Dest)
	err = uncompress(reader, config, mode >= modUPDATE)
	if err != nil {
		verboseFatal(err)
	}

	if mode >= modUPDATE {
		err = writeDestConfig(config)
		if err != nil {
			verboseFatal(err)
		}
	}

	verbosef("Running %s...", config.Run)
	err = run(config)
	if err != nil {
		verboseFatal(err)
	}
}
