// Package reader contains reader utilities for reading from a file
package reader

import (
	"io"
	"os"
)

type fileReader struct {
	rOs *os.File
}

// FileReader interface gives you some options for reading from a file
type FileReader interface {
	// ReadData func provides reading data from file by defining custom pos & seek option
	ReadData(pos int64, len int, seek int) ([]byte, error)
	// ReadAllData func provides reading all data from
	ReadAllData() ([]byte, error)
}

// NewFileReader func provides new instance of FileReader interface with unique memory addresses of its objects
func NewFileReader(filePath string) (FileReader, error) {
	var err error
	rOs := new(os.File)

	if rOs, err = os.OpenFile(filePath, os.O_RDONLY, 0444); err != nil {
		return &fileReader{}, err
	}

	return &fileReader{
		rOs: rOs,
	}, nil
}

// ReadData func provides reading data from file by defining custom pos & seek option
func (r *fileReader) ReadData(pos int64, len int, seek int) ([]byte, error) {
	buff := make([]byte, len)

	if _, err := r.rOs.Seek(pos, seek); err != nil {
		return nil, err
	}

	if _, err := r.rOs.Read(buff); err != nil {
		return nil, err
	}

	return buff, nil
}

func (r *fileReader) ReadAllData() ([]byte, error) {
	if buff, err := io.ReadAll(r.rOs); err != nil {
		return nil, err
	} else {
		return buff, nil
	}
}
