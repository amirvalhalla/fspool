package file

import (
	"io"
	"io/fs"
	"os"
)

// File override os.File interface of golang with RW interfaces
type File interface {
	io.Writer
	io.WriterAt
	io.WriteCloser
	io.WriteSeeker
	io.StringWriter
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Closer
	io.ReaderFrom
	io.ReadSeekCloser
	io.ReadSeeker
	Stat() (os.FileInfo, error)
	Sync() error
}

// FileInfo override fs.FileInfo interface of golang with RW interfaces
type FileInfo interface {
	fs.FileInfo
}

// FileHelper override static functions of os package
type FileHelper interface {
	Stat(path string) (os.FileInfo, error)
	IsNotExist(err error) bool
	MkdirAll(path string, perm os.FileMode) error
}
