// Package reader contains reader utilities for reading from a file
package reader

import (
	"errors"
	"os"
)

var (
	ErrCouldNotOpenFile = errors.New("could not open file")
	ErrCouldNotSeek     = errors.New("could not seek")
	ErrCouldNotRead     = errors.New("could not read data")
	ErrCouldNotReadAll  = errors.New("could not read all data")
	ErrCouldNotClose    = errors.New("could not close")
)

type fileReader struct {
	rOs *os.File
}

// FileReader interface gives you some options for reading from a file
type FileReader interface {
	// ReadData func provides reading data from file by defining custom pos & seek option
	ReadData(offset int64, len int, seek int) ([]byte, error)
	// ReadAllData func provides reading all data from file
	ReadAllData() ([]byte, error)
	// Close func provides close reader instance
	Close() error
}

// NewFileReader func provides new instance of FileReader interface with unique memory addresses of its objects
func NewFileReader(filePath string) (FileReader, error) {
	var err error
	rOs := new(os.File)

	if rOs, err = os.OpenFile(filePath, os.O_RDONLY, 0444); err != nil {
		return &fileReader{}, ErrCouldNotOpenFile
	}

	return &fileReader{
		rOs: rOs,
	}, nil
}

// ReadData func provides reading data from file by defining custom pos & seek option
func (r *fileReader) ReadData(offset int64, len int, seek int) ([]byte, error) {
	buff := make([]byte, len)

	if _, err := r.rOs.Seek(offset, seek); err != nil {
		return nil, ErrCouldNotSeek
	}

	if _, err := r.rOs.Read(buff); err != nil {
		return nil, ErrCouldNotRead
	}

	return buff, nil
}

// ReadAllData func provides reading all data from file
func (r *fileReader) ReadAllData() ([]byte, error) {
	var buffSize int64 = 0

	if fInfo, err := r.rOs.Stat(); err != nil {
		return nil, ErrCouldNotReadAll
	} else {
		buffSize = fInfo.Size()
	}

	buff := make([]byte, buffSize)
	if _, err := r.rOs.Read(buff); err != nil {
		return nil, ErrCouldNotReadAll
	}

	return buff, nil
}

// Close func provides close reader instance
func (r *fileReader) Close() error {
	if err := r.rOs.Close(); err != nil {
		return ErrCouldNotClose
	} else {
		return nil
	}
}
