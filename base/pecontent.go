package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"debug/pe"
	"errors"
	"io"
	"os"
)

func openExeContent(path string) (io.Closer, *peConfig, *tar.Reader, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, nil, nil, err
	}

	finfo, err := file.Stat()
	if err != nil {
		return nil, nil, nil, err
	}

	pefile, err := pe.NewFile(file)
	if err != nil {
		return nil, nil, nil, err
	}

	config, reader, err := getEmbedContent(file, pefile, finfo.Size())
	if err != nil {
		return nil, nil, nil, err
	}

	return file, config, reader, nil
}

func getEmbedContent(file *os.File, pefile *pe.File, size int64) (*peConfig, *tar.Reader, error) {
	var max int64
	for _, sec := range pefile.Sections {
		config, reader, err := readSection(file, sec.Open(), int64(sec.Offset), int64(sec.Size))
		if err == nil {
			return config, reader, nil
		} else if err.Error() != "ENOTCONFIG" {
			return nil, nil, err
		}

		// Otherwise move end of file pointer
		end := int64(sec.Offset + sec.Size)
		if end > max {
			max = end
		}
	}

	// No zip file within binary, try appended to end
	section := io.NewSectionReader(file, max, size-max)
	return readSection(file, section, max, section.Size())
}

func readSection(file *os.File, sec io.ReadSeeker, offset int64, size int64) (*peConfig, *tar.Reader, error) {
	raw := make([]byte, 11)
	_, err := sec.Read(raw)
	if err != nil {
		return nil, nil, err
	}

	buf := bytes.NewBuffer(raw)
	if !isConfigPart(buf.String()) {
		return nil, nil, errors.New("ENOTCONFIG")
	}

	_, err = sec.Seek(0, 0)
	if err != nil {
		return nil, nil, err
	}

	len, config, err := readConfig(sec)
	if err != nil {
		return nil, nil, err
	}

	zipSection := io.NewSectionReader(file, offset+int64(len), size-int64(len))
	gzipReader, err := gzip.NewReader(zipSection)
	tarReader := tar.NewReader(gzipReader)
	if err != nil {
		return nil, nil, err
	}

	return config, tarReader, nil
}
