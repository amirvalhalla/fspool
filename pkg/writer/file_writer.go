// Package writer contains writer utilities for writing into a file
package writer

import (
	"errors"
	"github.com/amirvalhalla/fspool/pkg/file"
	"github.com/google/uuid"
	"sync"
)

var (
	ErrFileWriterCouldNotSeek  = errors.New("package writer - could not seek")
	ErrFileWriterCouldNotWrite = errors.New("package writer - could not write data into file")
	ErrFileWriterCouldNotClose = errors.New("package writer - could not close")
	ErrFileWriterCouldNotSync  = errors.New("package writer - could not sync")
)

type fileWriter struct {
	id    uuid.UUID
	wFile file.File
	rwMu  sync.RWMutex
}

// FileWriter interface gives you some options for writing into a file
type FileWriter interface {
	// Write will write or update raw data into file
	Write(rawData []byte, offset int64, seek int) error
	// Sync will sync data from in-memory to disk
	Sync() error
	// GetId return id of FileWriter
	GetId() uuid.UUID
	// Close func provides close writer instance
	Close() error
}

// NewFileWriter func provides new instance of FileWriter interface with unique memory addresses of its objects
func NewFileWriter(file file.File) (FileWriter, uuid.UUID) {
	id := uuid.New()
	return &fileWriter{
		id:    id,
		wFile: file,
	}, id
}

// Write will write or update raw data into file
func (w *fileWriter) Write(rawData []byte, offset int64, seek int) error {
	w.rwMu.Lock()
	defer w.rwMu.Unlock()

	if _, err := w.wFile.Seek(offset, seek); err != nil {
		return ErrFileWriterCouldNotSeek
	}

	if _, err := w.wFile.Write(rawData); err != nil {
		return ErrFileWriterCouldNotWrite
	}

	return nil
}

// Sync will sync data from in-memory to disk
func (w *fileWriter) Sync() error {
	w.rwMu.Lock()
	defer w.rwMu.Unlock()

	if err := w.wFile.Sync(); err != nil {
		return ErrFileWriterCouldNotSync
	}

	return nil
}

// GetId return id of FileWriter
func (w *fileWriter) GetId() uuid.UUID {
	return w.id
}

// Close func provides close writer instance
func (w *fileWriter) Close() error {
	w.rwMu.Lock()
	defer w.rwMu.Unlock()

	if err := w.wFile.Close(); err != nil {
		return ErrFileWriterCouldNotClose
	}

	return nil
}
