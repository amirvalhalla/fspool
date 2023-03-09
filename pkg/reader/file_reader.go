// Package reader contains reader utilities for reading from a file
package reader

import (
	"errors"
	"github.com/amirvalhalla/fspool/pkg/file"
	"github.com/google/uuid"
	"sync"
)

var (
	ErrFileReaderCouldNotSeek        = errors.New("package reader - could not seek")
	ErrFileReaderCouldNotRead        = errors.New("package reader - could not read data")
	ErrFileReaderCouldNotReadAllData = errors.New("package reader - could not read all data")
	ErrFileReaderCouldNotClose       = errors.New("package reader - could not close")
	ErrFileReaderCouldNotGetFileStat = errors.New("package reader - could not get file stat")
)

type fileReader struct {
	id    uuid.UUID
	rFile file.File
	rwMu  sync.RWMutex
}

// FileReader interface gives you some options for reading from a file
type FileReader interface {
	// ReadData func provides reading data from file by defining custom pos & seek option
	ReadData(offset int64, len int, seek int) ([]byte, error)
	// ReadAllData func provides reading all data from file
	ReadAllData() ([]byte, error)
	// GetId return id of FileReader
	GetId() uuid.UUID
	// Close func provides close reader instance
	Close() error
}

// NewFileReader func provides new instance of FileReader interface with unique memory addresses of its objects
func NewFileReader(file file.File) (FileReader, uuid.UUID) {
	id := uuid.New()

	return &fileReader{
		id:    id,
		rFile: file,
	}, id
}

// ReadData func provides reading data from file by defining custom pos & seek option
func (r *fileReader) ReadData(offset int64, len int, seek int) ([]byte, error) {
	r.rwMu.RLock()
	defer r.rwMu.RUnlock()

	buff := make([]byte, len)

	if _, err := r.rFile.Seek(offset, seek); err != nil {
		return nil, ErrFileReaderCouldNotSeek
	}

	if _, err := r.rFile.Read(buff); err != nil {
		return nil, ErrFileReaderCouldNotRead
	}

	return buff, nil
}

// ReadAllData func provides reading all data from file
func (r *fileReader) ReadAllData() ([]byte, error) {
	r.rwMu.RLock()
	defer r.rwMu.RUnlock()

	var buffSize int64 = 0

	if fInfo, err := r.rFile.Stat(); err != nil {
		return nil, ErrFileReaderCouldNotGetFileStat
	} else {
		buffSize = fInfo.Size()
	}

	buff := make([]byte, buffSize)
	if _, err := r.rFile.Read(buff); err != nil {
		return nil, ErrFileReaderCouldNotReadAllData
	}

	return buff, nil
}

// GetId return id of FileReader
func (r *fileReader) GetId() uuid.UUID {
	return r.id
}

// Close func provides close reader instance
func (r *fileReader) Close() error {
	r.rwMu.RLock()
	defer r.rwMu.RUnlock()

	if err := r.rFile.Close(); err != nil {
		return ErrFileReaderCouldNotClose
	}

	return nil
}
