package vfs

import (
	"io"
	"os"
)

type Filesystem interface {
	Create(filename string) (File, error)
	Open(filename string) (File, error)
	OpenFile(filename string, flag int, perm os.FileMode) (File, error)
	Stat(filename string) (os.FileInfo, error)
	Rename(oldpath, newpath string) error
	Remove(filename string) error
}

type File interface {
	Name() string
	Fd() uintptr
	io.Writer
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Closer
}
