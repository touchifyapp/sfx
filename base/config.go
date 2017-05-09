package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type peConfig struct {
	ID      string
	Version string
	Dest    string
	Run     string
	Args    string
}

func isConfigPart(part string) bool {
	return strings.HasPrefix(part, "[sfxconfig]")
}

func isConfigPartEnd(part string) bool {
	return strings.HasPrefix(part, "[/sfxconfig]")
}

func readDestConfig(config *peConfig) (*peConfig, error) {
	path := filepath.Join(config.Dest, "SFXCONFIG")
	file, err := os.OpenFile(path, os.O_RDONLY, 0400)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, err
	}

	_, destConfig, err := readConfig(file)
	if err != nil {
		return nil, err
	}

	err = file.Close()
	if err != nil {
		return nil, err
	}

	return destConfig, nil
}

func writeDestConfig(config *peConfig) error {
	conf := serializeConfig(config)

	path := filepath.Join(config.Dest, "SFXCONFIG")
	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	file.WriteString(conf)

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

func readConfig(reader io.ReadSeeker) (int, *peConfig, error) {
	len, config, err := parseConfig(reader)
	if err != nil {
		return len, nil, err
	}

	if config.Dest == "" {
		config.Dest = filepath.Join(os.TempDir(), config.ID)
	}

	return len, config, nil
}

func parseConfig(reader io.ReadSeeker) (int, *peConfig, error) {
	config := new(peConfig)
	size := 12

	// Skip config start tag
	_, err := reader.Seek(12, 0) // [sfxconfig] + \n
	if err != nil {
		return size, nil, err
	}

	len, name, value, err := parseNext(reader)
	for err == nil {
		size = size + len
		setValue(config, name, value)

		len, name, value, err = parseNext(reader)
	}

	if err != nil {
		if err.Error() == "EINVALIDCONF" {
			// Restore reader to position before error was thrown
			_, err := reader.Seek(-1, 1)
			if err != nil {
				return size, nil, err
			}

			isEnd, err := isPartEnd(reader)
			if err != nil {
				return size, nil, err
			}

			if isEnd {
				size = size + 12
			} else {
				return size, nil, errors.New("EINVALIDCONF")
			}
		} else {
			return size, nil, err
		}
	}

	verbosef("Config found: %+v (%d)", config, size)
	return size, config, nil
}

func serializeConfig(config *peConfig) string {
	conf := "[sfxconfig]\n" +
		"ID=" + config.ID + "\n" +
		"Run=" + config.Run + "\n" +
		"Version=" + config.Version + "\n"

	if config.Dest != "" {
		conf = conf + "Dest=" + config.Dest + "\n"
	}

	if config.Args != "" {
		conf = conf + "Args=" + config.Args + "\n"
	}

	conf = conf + "[/sfxconfig]"

	return conf
}

func setValue(config *peConfig, name string, value string) {
	switch name {
	case "ID":
		config.ID = value

	case "Version":
		config.Version = value

	case "Dest":
		config.Dest = value

	case "Run":
		config.Run = value

	case "Args":
		config.Args = value
	}
}

func parseNext(reader io.ReadSeeker) (int, string, string, error) {
	size := 0

	len, name, err := parseName(reader)
	if err != nil {
		return size, "", "", err
	}

	size = size + len

	len, value, err := parseValue(reader)
	if err != nil {
		return size, "", "", err
	}

	size = size + len

	return size, name, value, nil
}

func parseName(reader io.ReadSeeker) (int, string, error) {
	name := ""
	char := ""
	len := 0
	raw := make([]byte, 1)

	for char != "=" {
		name = name + char

		_, err := reader.Read(raw)
		if err != nil {
			return len, "", err
		}

		len = len + 1
		char = bytes.NewBuffer(raw).String()
		if char == "\n" || char == "[" {
			return len, "", errors.New("EINVALIDCONF")
		}
	}

	return len, name, nil
}

func parseValue(reader io.ReadSeeker) (int, string, error) {
	value := ""
	char := ""
	len := 0
	raw := make([]byte, 1)

	for char != "\n" {
		value = value + char

		_, err := reader.Read(raw)
		if err != nil {
			return len, "", err
		}

		len = len + 1
		char = bytes.NewBuffer(raw).String()
		if char == "=" || char == "[" {
			return len, "", errors.New("EINVALIDCONF")
		}
	}

	return len, value, nil
}

func isPartEnd(reader io.ReadSeeker) (bool, error) {
	raw := make([]byte, 12)
	_, err := reader.Read(raw)
	if err != nil {
		return false, err
	}

	return isConfigPartEnd(bytes.NewBuffer(raw).String()), nil
}
