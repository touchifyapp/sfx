package main

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func appendTarFile(exefile *os.File) error {
	tmpTarPath := filepath.Join(os.TempDir(), fmt.Sprintf("sfx-%d.tar.gz", time.Now().Unix()))
	tmpTarfile, err := os.Create(tmpTarPath)
	if err != nil {
		return err
	}

	defer func() {
		tmpTarfile.Close()
		os.Remove(tmpTarPath)
	}()

	verbosef("Compress package using level %d", args.Compress)
	gzipWriter, err := gzip.NewWriterLevel(tmpTarfile, args.Compress)
	if err != nil {
		return err
	}

	tarWriter := tar.NewWriter(gzipWriter)

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

		tarFileName := strings.TrimPrefix(path, dirAbs+string(filepath.Separator))

		tarFileHeader, err := tar.FileInfoHeader(info, tarFileName)
		if err != nil {
			return err
		}

		// tarFileHeader.Name = tarFileName

		err = tarWriter.WriteHeader(tarFileHeader)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}

		_, err = io.Copy(tarWriter, srcFile)
		if err != nil {
			return err
		}

		srcFile.Close()

		return nil
	})

	if err != nil {
		return err
	}

	err = tarWriter.Close()
	if err != nil {
		return err
	}

	err = gzipWriter.Close()
	if err != nil {
		return err
	}

	err = tmpTarfile.Sync()
	if err != nil {
		return err
	}

	_, err = tmpTarfile.Seek(0, 0)
	if err != nil {
		return err
	}

	_, err = exefile.Seek(0, 2)
	if err != nil {
		return err
	}

	_, err = io.Copy(exefile, tmpTarfile)
	if err != nil {
		return err
	}

	return nil
}
