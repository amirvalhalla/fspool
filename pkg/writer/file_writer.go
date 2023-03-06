// Package writer contains writer utilities for writing into a file
package writer

import (
	"errors"
	"github.com/amirvalhalla/fspool/pkg/file"
	"sync"
)

var (
	ErrFileWriterCouldNotSeek  = errors.New("package writer - could not seek")
	ErrFileWriterCouldNotWrite = errors.New("package writer - could not write data into file")
	ErrFileWriterCouldNotClose = errors.New("package writer - could not close")
)

type fileWriter struct {
	wFile file.File
	rwMu  sync.RWMutex
}

// FileWriter interface gives you some options for writing into a file
type FileWriter interface {
	// Write will write or update raw data into file
	Write(rawData []byte, offset int64, seek int) error
	// Close func provides close writer instance
	Close() error
}

// NewFileWriter func provides new instance of FileWriter interface with unique memory addresses of its objects
func NewFileWriter(file file.File) FileWriter {
	return &fileWriter{
		wFile: file,
	}
}

// Write will write or update raw data into file
func (w *fileWriter) Write(rawData []byte, offset int64, seek int) error {
	w.rwMu.Lock()
	defer w.rwMu.Unlock()

	_, err := w.wFile.Seek(offset, seek)

	if err != nil {
		return ErrFileWriterCouldNotSeek
	}

	_, err = w.wFile.Write(rawData)

	if err != nil {
		return ErrFileWriterCouldNotWrite
	}

	return nil
}

// Close func provides close writer instance
func (w *fileWriter) Close() error {
	w.rwMu.Lock()
	defer w.rwMu.Unlock()

	if err := w.wFile.Close(); err != nil {
		return ErrFileWriterCouldNotClose
	} else {
		return nil
	}
}
