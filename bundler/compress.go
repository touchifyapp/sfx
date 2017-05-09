package main

import (
	"archive/zip"
	"compress/flate"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func appendZipFile(exefile *os.File) error {
	tmpZipPath := filepath.Join(os.TempDir(), fmt.Sprintf("sfx-%d.zip", time.Now().Unix()))
	tmpZipfile, err := os.Create(tmpZipPath)
	if err != nil {
		return err
	}

	defer func() {
		tmpZipfile.Close()
		os.Remove(tmpZipPath)
	}()

	zipWriter := zip.NewWriter(tmpZipfile)
	if args.Compress > 0 {
		verbosef("Compress package using level %d", args.Compress)
		zipWriter.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
			return flate.NewWriter(out, flate.BestCompression)
		})
	}

	dirAbs, err := filepath.Abs(args.Dir)
	if err != nil {
		return err
	}

	err = filepath.Walk(dirAbs, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return errors.New(dirAbs + " does not exist on this computer")
		}

		if path == dirAbs {
			return nil
		}

		zipFileName := strings.TrimPrefix(path, dirAbs+string(filepath.Separator))

		if info.IsDir() {

			return nil
		}

		zipFileHeader, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		zipFileHeader.Name = zipFileName

		zipFileWriter, err := zipWriter.CreateHeader(zipFileHeader)
		if err != nil {
			return err
		}

		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}

		_, err = io.Copy(zipFileWriter, srcFile)
		if err != nil {
			return err
		}

		srcFile.Close()

		return nil
	})

	if err != nil {
		return err
	}

	err = zipWriter.Close()
	if err != nil {
		return err
	}

	err = tmpZipfile.Sync()
	if err != nil {
		return err
	}

	_, err = tmpZipfile.Seek(0, 0)
	if err != nil {
		return err
	}

	_, err = exefile.Seek(0, 2)
	if err != nil {
		return err
	}

	_, err = io.Copy(exefile, tmpZipfile)
	if err != nil {
		return err
	}

	return nil
}
