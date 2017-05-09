package main

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
)

func uncompress(reader *tar.Reader, config *peConfig, force bool) error {
	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		dest := filepath.Join(config.Dest, header.Name)
		info := header.FileInfo()

		if info.IsDir() {
			err = os.MkdirAll(dest, info.Mode())
			if err != nil {
				return err
			}

			continue
		}

		if !force {
			destInfo, err := os.Stat(dest)
			if err != nil && !os.IsNotExist(err) {
				return err
			}

			if destInfo.ModTime().After(info.ModTime()) {
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
