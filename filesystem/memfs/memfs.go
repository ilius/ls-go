package memfs

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/ilius/ls-go/filesystem/paths"
)

const separator = filepath.Separator

var (
	ErrNotLink = errors.New("not a link")
	ErrIsDir   = errors.New("is a directory")
)

// Memory a very convenient filesystem based on memory files
type Memory struct {
	*paths.LocalFilePath

	s *storage

	// working directory
	wd string

	// tempCount int
}

// New returns a new Memory filesystem.
func New() *Memory {
	fs := &Memory{
		LocalFilePath: mypath,

		s: newStorage(),
	}
	return fs
}

// AddFile creates a file and writes the given data
func (m *Memory) AddFile(path string, mode fs.FileMode, data []byte) error {
	f, err := m.s.New(path, mode, os.O_CREATE)
	if err != nil {
		return err
	}
	f.content.bytes = data
	return nil
}

// MkdirAll create a new directory and its parents if needed
func (m *Memory) MkdirAll(path string, perm fs.FileMode) error {
	_, err := m.s.New(path, perm|os.ModeDir, 0)
	return err
}

func (m *Memory) Symlink(target string, path string, perm fs.FileMode) error {
	_, err := m.Stat(path)
	if err == nil {
		return os.ErrExist
	}

	if !os.IsNotExist(err) {
		return err
	}

	return m.AddFile(
		path,
		perm&os.ModeSymlink,
		[]byte(target),
	)
}

// Open opens the named object for reading.
func (m *Memory) Open(name string) (fs.File, error) {
	f, ok := m.s.Get(m.Join(m.wd, name))
	if !ok {
		return nil, &fs.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
	}
	return f, nil
}

// Stat returns a FileInfo for the given name.
func (m *Memory) Stat(name string) (fs.FileInfo, error) {
	f, ok := m.s.Get(m.Join(m.wd, name))
	if !ok {
		return nil, &fs.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
	}
	return f.Stat()
}

// ReadFile reads the whole object into a byte slice.
func (m *Memory) ReadFile(name string) ([]byte, error) {
	f, ok := m.s.Get(m.Join(m.wd, name))
	if !ok {
		return nil, &fs.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
	}
	if f.mode.IsDir() {
		return nil, &fs.PathError{Op: "open", Path: name, Err: ErrIsDir}
	}
	return f.content.bytes, nil
}

func (m *Memory) resolveLink(fullpath string, f *file) (target string, isLink bool) {
	if !isSymlink(f.mode) {
		return fullpath, false
	}

	target = string(f.content.bytes)
	if !m.IsAbs(target) {
		target = m.Join(filepath.Dir(fullpath), target)
	}

	return target, true
}

type DirEntry struct {
	fi fs.FileInfo
}

// Name returns the name of the file (or subdirectory) described by the entry.
// This name is only the final element of the path (the base name), not the entire path.
// For example, Name would return "hello.go" not "home/gopher/hello.go".
func (e *DirEntry) Name() string {
	return e.fi.Name()
}

// IsDir reports whether the entry describes a directory.
func (e *DirEntry) IsDir() bool {
	return e.fi.IsDir()
}

// Type returns the type bits for the entry.
// The type bits are a subset of the usual FileMode bits, those returned by the FileMode.Type method.
func (e *DirEntry) Type() fs.FileMode {
	return e.fi.Mode()
}

// Info returns the FileInfo for the file or subdirectory described by the entry.
// The returned FileInfo may be from the time of the original directory read
// or from the time of the call to Info. If the file has been removed or renamed
// since the directory read, Info may return an error satisfying errors.Is(err, ErrNotExist).
// If the entry denotes a symbolic link, Info reports the information about the link itself,
// not the link's target.
func (e *DirEntry) Info() (fs.FileInfo, error) {
	return e.fi, nil
}

// ReadDir reads the directory and returns a list of DirEntry.
func (m *Memory) ReadDir(name string) ([]fs.DirEntry, error) {
	if f, has := m.s.Get(name); has {
		if target, isLink := m.resolveLink(name, f); isLink {
			return m.ReadDir(target)
		}
	}

	children := m.s.Children(name)
	entries := make([]fs.DirEntry, len(children))
	for index, file := range children {
		fi, _ := file.Stat()
		entries[index] = &DirEntry{fi}
	}

	return entries, nil
}

// ReadLink returns the destination of the named symbolic link. If there is an error, it will be of type *os.PathError.
func (m *Memory) ReadLink(name string) (string, error) {
	f, ok := m.s.Get(name)
	if !ok {
		return "", &fs.PathError{Op: "readLink", Path: name, Err: os.ErrNotExist}
	}
	target, isLink := m.resolveLink(name, f)
	if !isLink {
		return "", &fs.PathError{Op: "readLink", Path: name, Err: ErrNotLink}
	}
	return target, nil
}

// WorkDir returns the path of working directory
func (fs *Memory) WorkDir() string {
	return fs.wd
}

// UserHomeDir returns the current user's home directory.
func (fs *Memory) UserHomeDir() (string, error) {
	return string(separator), nil
}

// CountDirContents: returns the number of files/directories direnctly under a given directory
func (fs *Memory) CountDirContents(_ string) (int, error) {
	// FIXME: it's counting the root dir, not the given dir as argument
	return len(fs.s.children), nil
}
