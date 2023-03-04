// Package reader contains reader utilities for reading from a file
package reader

import (
	"errors"
	"io"
	"io/fs"
	"os"
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
	rFile File
	rwMu  sync.RWMutex
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

// File override os.File interface of golang with ROnly interfaces
type File interface {
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Closer
	io.ByteReader
	io.ReaderFrom
	io.RuneReader
	io.RuneScanner
	io.ReadSeekCloser
	io.ReadSeeker
	fs.FileInfo
	Stat() (os.FileInfo, error)
}

// NewFileReader func provides new instance of FileReader interface with unique memory addresses of its objects
func NewFileReader(file File) FileReader {
	return &fileReader{
		rFile: file,
	}
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

// Close func provides close reader instance
func (r *fileReader) Close() error {
	r.rwMu.RLock()
	defer r.rwMu.RUnlock()

	if err := r.rFile.Close(); err != nil {
		return ErrFileReaderCouldNotClose
	} else {
		return nil
	}
}
