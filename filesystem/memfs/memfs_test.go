package memfs

import (
	"io/fs"
	"testing"
	"time"

	"github.com/ilius/is/v2"
	"github.com/ilius/ls-go/iface"
)

func init() {
	var _ iface.FileSystem = New()
}

func TestAddRead(t *testing.T) {
	is := is.New(t)
	m := New()
	now := time.Now()
	is.NotErr(
		m.AddFile("a.txt", fs.FileMode(0o644), []byte(`hello`)),
	)
	info, err := m.Stat("a.txt")
	is.NotErr(err)
	is.Equal("a.txt", info.Name())
	is.False(info.IsDir())
	is.Equal(5, info.Size())
	is.Equal(info.Mode(), fs.FileMode(0o644))
	tm := info.ModTime()
	is.True(tm.Sub(now) < time.Second)
	f, err := m.Open("a.txt")
	is.NotErr(err)
	is.NotNil(f)
	info2, err := f.Stat()
	is.NotErr(err)
	is.Equal(info, info2)
	var data []byte
	n, err := f.Read(data)
	is.ErrMsg(err, "read not supported")
	is.Equal(0, n)
	// is.NotErr(err)
	// is.Equal(5, n)
	// is.Equal([]byte(`hello`), data)
}

func TestLink(t *testing.T) {
	is := is.New(t)
	m := New()
	is.NotErr(
		m.AddFile("a.txt", fs.FileMode(0o644), []byte(`hello`)),
	)
	is.NotErr(
		m.Symlink("a.txt", "b.txt", fs.FileMode(0o644)),
	)
}
