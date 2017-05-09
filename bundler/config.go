package main

import (
	"bytes"
	"io"
	"os"
)

func appendConfigFile(exefile *os.File) error {
	config := "[sfxconfig]\n" +
		"ID=" + args.ID + "\n" +
		"Run=" + args.Run + "\n" +
		"Version=" + args.Version + "\n"

	if args.Dest != "" {
		config = config + "Dest=" + args.Dest + "\n"
	}

	if args.Args != "" {
		config = config + "Args=" + args.Args + "\n"
	}

	config = config + "[/sfxconfig]"

	buf := bytes.NewBufferString(config)
	reader := bytes.NewReader(buf.Bytes())

	_, err := exefile.Seek(0, 2)
	if err != nil {
		return err
	}

	_, err = io.Copy(exefile, reader)
	if err != nil {
		return err
	}

	return nil
}
