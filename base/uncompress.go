package main

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func uncompress(reader *zip.Reader, config *peConfig) error {
	for _, file := range reader.File {
		dest := filepath.Join(config.Dest, file.Name)
		err := copyFile(reader, file, dest)
		if err != nil {
			return err
		}
	}

	return nil
}

func copyFile(reader *zip.Reader, file *zip.File, dest string) error {
	fileReader, err := file.Open()
	if err != nil {
		return err
	}

	defer fileReader.Close()

	err = os.MkdirAll(filepath.Dir(dest), 0777)
	if err != nil {
		return err
	}

	destFile, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, 0700)
	if err != nil {
		return err
	}

	defer destFile.Close()

	_, err = io.Copy(destFile, fileReader)
	if err != nil {
		return err
	}

	return nil
}
