package iface

import (
	"io/fs"

	"github.com/ilius/ls-go/common"
)

// implements the fs.FS interface.
type FileSystem interface {
	// Open opens the named object for reading.
	Open(name string) (fs.File, error)

	// Stat returns a FileInfo for the given name.
	Stat(name string) (fs.FileInfo, error)

	// ReadDir reads the directory and returns a list of DirEntry.
	ReadDir(name string) ([]fs.DirEntry, error)

	// ReadLink returns the destination of the named symbolic link. If there is an error, it will be of type *os.PathError.
	ReadLink(name string) (string, error)

	// WorkDir returns the path of working directory
	WorkDir() string

	// UserHomeDir returns the current user's home directory.
	UserHomeDir() (string, error)

	// CountDirContents: returns the number of files/directories direnctly under a given directory
	CountDirContents(name string) (int, error)

	// abstraction of "path/filepath" methods

	// same as filepath.Abs
	Abs(path string) (string, error)

	// same as filepath.Dir, but also work with path ending with a slash/sep
	Dir(path string) string

	// same as filepath.IsAbs
	IsAbs(path string) bool

	// same as filepath.Join
	Join(elem ...string) string

	// same as filepath.Rel
	Rel(basepath, targpath string) (string, error)

	SplitAll(path string) []string

	// same as Join, but give a color to path sep (slash or backslash)
	JoinColor(color string, reset string, elem ...string) string

	SplitExt(filename string) *common.ParsedName
}
