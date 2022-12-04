package filesystem

import (
	"io/fs"
	"os"

	"github.com/ilius/ls-go/filesystem/paths"
)

func NewLocalFileSystem() *LocalFileSystem {
	return &LocalFileSystem{}
}

// FS implements the fs.FS interface.
type LocalFileSystem struct {
	paths.LocalFilePath
}

// Open opens the named object for reading.
func (*LocalFileSystem) Open(name string) (fs.File, error) {
	return os.Open(name)
}

// Stat returns a FileInfo for the given name.
func (*LocalFileSystem) Stat(name string) (fs.FileInfo, error) {
	return os.Stat(name)
}

// ReadDir reads the directory and returns a list of DirEntry.
func (*LocalFileSystem) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(name)
}

// ReadLink returns the destination of the named symbolic link. If there is an error, it will be of type *os.PathError.
func (*LocalFileSystem) ReadLink(name string) (string, error) {
	return os.Readlink(name)
}

// WorkDir returns the path of working directory
func (*LocalFileSystem) WorkDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return cwd
}

// UserHomeDir returns the current user's home directory.
func (*LocalFileSystem) UserHomeDir() (string, error) {
	return os.UserHomeDir()
}

// CountDirContents: returns the number of files/directories direnctly under a given directory
func (*LocalFileSystem) CountDirContents(name string) (int, error) {
	file, err := os.Open(name)
	if err != nil {
		return 0, err
	}
	list, err := file.Readdirnames(-1)
	file.Close()
	if err != nil {
		return 0, err
	}
	return len(list), nil
}
