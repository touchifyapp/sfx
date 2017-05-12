package main

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
	"time"
)

func uncompress(reader *tar.Reader, config *peConfig, newModTime time.Time) error {
	err := os.MkdirAll(config.Dest, 0777)
	if err != nil {
		return err
	}

	checkDates := newModTime.After(time.Unix(0, 0))

	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		dest := filepath.Join(config.Dest, header.Linkname)
		info := header.FileInfo()

		if info.IsDir() {
			err = os.MkdirAll(dest, info.Mode())
			if err != nil {
				return err
			}

			continue
		}

		if checkDates {
			destInfo, err := os.Stat(dest)
			if err != nil && !os.IsNotExist(err) {
				return err
			}

			if destInfo.ModTime().After(newModTime) {
				return nil
			}
		}

		destFile, err := os.OpenFile(dest, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}

		defer destFile.Close()

		_, err = io.Copy(destFile, reader)
		if err != nil {
			return err
		}
	}

	return nil
}
