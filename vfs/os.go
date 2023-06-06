package vfs

import "os"

var _ Filesystem = &osfs{}

type osfs struct {
}

func NewOSFS() Filesystem {
	return &osfs{}
}

func (fs *osfs) Create(filename string) (File, error) {
	return os.Create(filename)
}
func (fs *osfs) Open(filename string) (File, error) {
	return os.Open(filename)
}
func (fs *osfs) OpenFile(filename string, flag int, perm os.FileMode) (File, error) {
	return os.OpenFile(filename, flag, perm)
}
func (fs *osfs) Stat(filename string) (os.FileInfo, error) {
	return os.Stat(filename)
}
func (fs *osfs) Rename(oldpath, newpath string) error {
	return os.Rename(oldpath, newpath)
}
func (fs *osfs) Remove(filename string) error {
	return os.Remove(filename)
}
