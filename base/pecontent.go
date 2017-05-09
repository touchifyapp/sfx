package main

import (
	"archive/zip"
	"bytes"
	"debug/pe"
	"errors"
	"io"
	"os"
)

func openExeContent(path string) (io.Closer, *peConfig, *zip.Reader, error) {
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

func getEmbedContent(file *os.File, pefile *pe.File, size int64) (*peConfig, *zip.Reader, error) {
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

func readSection(file *os.File, sec io.ReadSeeker, offset int64, size int64) (*peConfig, *zip.Reader, error) {
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
	zipfile, err := zip.NewReader(zipSection, zipSection.Size())
	if err != nil {
		return nil, nil, err
	}

	return config, zipfile, nil
}

// zipExeReaderPe treats the file as a Portable Exectuable binary
// (Windows executable) and attempts to find a zip archive.
func zipExeReaderPe(rda io.ReaderAt, size int64) (*zip.Reader, error) {
	file, err := pe.NewFile(rda)
	if err != nil {
		return nil, err
	}

	var max int64
	for _, sec := range file.Sections {
		// Check if this section has a zip file
		if zfile, err := zip.NewReader(sec, int64(sec.Size)); err == nil {
			return zfile, nil
		}

		// Otherwise move end of file pointer
		end := int64(sec.Offset + sec.Size)
		if end > max {
			max = end
		}
	}

	// No zip file within binary, try appended to end
	section := io.NewSectionReader(rda, max, size-max)
	return zip.NewReader(section, section.Size())
}
