package json

import (
	"testing"

	"github.com/ilius/is/v2"
)

func TestParseFileInfo(t *testing.T) {
	is := is.New(t)
	jstr := `{"mode":"-rw-r--r--","hard_links":1,"owner":"ilius","group":"users",
"size":8451,"mtime":"2022-11-29 10:49:32.639995101 +0330","name":"README.md","link_target":""}`
	info, err := ParseFileInfo([]byte(jstr))
	if !is.NotErr(err) {
		return
	}
	is.Equal("README.md", info.Name())
	is.Equal(0o644, info.Mode())
	{
		n, err := info.NumberOfHardLinks()
		is.NotErr(err)
		is.Equal(1, n)
	}
	is.Equal("ilius", info.Owner())
	is.Equal("users", info.Group())
	is.Equal(8451, info.Size())
	is.Equal("2022-11-29 10:49:32.639995101 +0330", info.ModTime().Format(timeFmt))
	is.Equal("2022-11-29 10:49:32.639995101 +0330", info.Time("mtime").Format(timeFmt))
	is.Nil(info.CTime())
	is.Nil(info.ATime())

	// is.Equal("", info.)
}
