package memfs

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/ilius/ls-go/filesystem/paths"
)

var mypath = &paths.LocalFilePath{}

type storage struct {
	files    map[string]*file
	children map[string]map[string]*file
}

func newStorage() *storage {
	return &storage{
		files:    make(map[string]*file, 0),
		children: make(map[string]map[string]*file, 0),
	}
}

func (s *storage) Has(path string) bool {
	path = clean(path)

	_, ok := s.files[path]
	return ok
}

func (s *storage) New(path string, mode fs.FileMode, flag int) (*file, error) {
	path = clean(path)
	if s.Has(path) {
		if !s.MustGet(path).mode.IsDir() {
			return nil, &os.PathError{
				Op:   "create",
				Path: path,
				Err:  os.ErrExist,
			}
		}

		return nil, nil
	}

	name := mypath.Base(path)

	f := &file{
		name:    name,
		content: &content{name: name},
		mode:    mode,
		flag:    flag,
	}

	s.files[path] = f
	s.createParent(path, mode, f)
	return f, nil
}

func (s *storage) createParent(path string, mode fs.FileMode, f *file) error {
	base := mypath.Dir(path)
	base = clean(base)
	if f.Name() == string(separator) {
		return nil
	}

	if _, err := s.New(base, mode.Perm()|os.ModeDir, 0); err != nil {
		return err
	}

	if _, ok := s.children[base]; !ok {
		s.children[base] = make(map[string]*file, 0)
	}

	s.children[base][f.Name()] = f
	return nil
}

func (s *storage) Children(path string) []*file {
	path = clean(path)

	l := make([]*file, 0)
	for _, f := range s.children[path] {
		l = append(l, f)
	}

	return l
}

func (s *storage) MustGet(path string) *file {
	f, ok := s.Get(path)
	if !ok {
		panic(fmt.Errorf("couldn't find %q", path))
	}

	return f
}

func (s *storage) Get(path string) (*file, bool) {
	path = clean(path)
	if !s.Has(path) {
		return nil, false
	}

	file, ok := s.files[path]
	return file, ok
}

func (s *storage) Remove(path string) error {
	path = clean(path)

	f, has := s.Get(path)
	if !has {
		return os.ErrNotExist
	}

	if f.mode.IsDir() && len(s.children[path]) != 0 {
		return fmt.Errorf("dir: %s contains files", path)
	}

	base, file := mypath.Split(path)
	base = mypath.Clean(base)

	delete(s.children[base], file)
	delete(s.files, path)
	return nil
}

func clean(path string) string {
	return mypath.Clean(mypath.FromSlash(path))
}

type content struct {
	name  string
	bytes []byte
}

func (c *content) WriteAt(p []byte, off int64) (int, error) {
	if off < 0 {
		return 0, &os.PathError{
			Op:   "writeat",
			Path: c.name,
			Err:  errors.New("negative offset"),
		}
	}

	prev := len(c.bytes)

	diff := int(off) - prev
	if diff > 0 {
		c.bytes = append(c.bytes, make([]byte, diff)...)
	}

	c.bytes = append(c.bytes[:off], p...)
	if len(c.bytes) < prev {
		c.bytes = c.bytes[:prev]
	}

	return len(p), nil
}

func (c *content) ReadAt(b []byte, off int64) (n int, err error) {
	if off < 0 {
		return 0, &os.PathError{
			Op:   "readat",
			Path: c.name,
			Err:  errors.New("negative offset"),
		}
	}

	size := int64(len(c.bytes))
	if off >= size {
		return 0, io.EOF
	}

	l := int64(len(b))
	if off+l > size {
		l = size - off
	}

	btr := c.bytes[off : off+l]
	if len(btr) < len(b) {
		err = io.EOF
	}
	n = copy(b, btr)

	return
}
