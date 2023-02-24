// Package writer contains writer utilities for writing into a file
package writer

import (
	"errors"
	"os"
)

var (
	ErrFileWriterCouldNotOpenFile = errors.New("could not open file")
	ErrFileWriterCouldNotSeek     = errors.New("could not seek")
	ErrFileWriterCouldNotWrite    = errors.New("could not write data into file")
	ErrFileWriterCouldNotClose    = errors.New("could not close")
)

type fileWriter struct {
	wOs *os.File
}

// FileWriter interface gives you some options for writing into a file
type FileWriter interface {
	// AddOrUpdateData will add or update raw data into file
	AddOrUpdateData(rawData []byte, offset int64, seek int) error
	// Close func provides close writer instance
	Close() error
}

// NewFileWriter func provides new instance of FileWriter interface with unique memory addresses of its objects
func NewFileWriter(filePath string) (FileWriter, error) {
	var err error
	wOs := new(os.File)

	if wOs, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666); err != nil {
		return &fileWriter{}, ErrFileWriterCouldNotOpenFile
	}

	return &fileWriter{
		wOs: wOs,
	}, nil
}

// AddOrUpdateData will add or update raw data into file
func (w *fileWriter) AddOrUpdateData(rawData []byte, offset int64, seek int) error {
	_, err := w.wOs.Seek(offset, seek)

	if err != nil {
		return ErrFileWriterCouldNotSeek
	}

	_, err = w.wOs.Write(rawData)

	if err != nil {
		return ErrFileWriterCouldNotWrite
	}

	return nil
}

// Close func provides close writer instance
func (w *fileWriter) Close() error {
	if err := w.wOs.Close(); err != nil {
		return ErrFileWriterCouldNotClose
	} else {
		return nil
	}
}
