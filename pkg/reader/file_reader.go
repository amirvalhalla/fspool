// Package reader contains reader utilities for reading from a file
package reader

import (
	"errors"
	"os"
)

var (
	ErrFileReaderCouldNotOpenFile = errors.New("package reader - could not open file")
	ErrFileReaderCouldNotSeek     = errors.New("package reader - could not seek")
	ErrFileReaderCouldNotRead     = errors.New("package reader - could not read data")
	ErrFileReaderCouldNotReadAll  = errors.New("package reader - could not read all data")
	ErrFileReaderCouldNotClose    = errors.New("package reader - could not close")
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
		return &fileReader{}, ErrFileReaderCouldNotOpenFile
	}

	return &fileReader{
		rOs: rOs,
	}, nil
}

// ReadData func provides reading data from file by defining custom pos & seek option
func (r *fileReader) ReadData(offset int64, len int, seek int) ([]byte, error) {
	buff := make([]byte, len)

	if _, err := r.rOs.Seek(offset, seek); err != nil {
		return nil, ErrFileReaderCouldNotSeek
	}

	if _, err := r.rOs.Read(buff); err != nil {
		return nil, ErrFileReaderCouldNotRead
	}

	return buff, nil
}

// ReadAllData func provides reading all data from file
func (r *fileReader) ReadAllData() ([]byte, error) {
	var buffSize int64 = 0

	if fInfo, err := r.rOs.Stat(); err != nil {
		return nil, ErrFileReaderCouldNotReadAll
	} else {
		buffSize = fInfo.Size()
	}

	buff := make([]byte, buffSize)
	if _, err := r.rOs.Read(buff); err != nil {
		return nil, ErrFileReaderCouldNotReadAll
	}

	return buff, nil
}

// Close func provides close reader instance
func (r *fileReader) Close() error {
	if err := r.rOs.Close(); err != nil {
		return ErrFileReaderCouldNotClose
	} else {
		return nil
	}
}
