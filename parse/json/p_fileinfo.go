package json

import (
	"io/fs"
	"strconv"
	"strings"
	"time"

	"github.com/ilius/ls-go/filesystem/paths"
	"github.com/ilius/ls-go/lsplatform"
	"github.com/ilius/ls-go/parse"
)

const timeFmt = "2006-01-02 15:04:05.999999999 Z0700"

var (
	platform = lsplatform.New()
	fspaths  = &paths.LocalFilePath{}
)

func NewFakeFileInfo() *FakeFileInfo {
	return &FakeFileInfo{
		sys: platform.EmptyFileInfoSys(),
	}
}

type FakeFileInfo struct {
	sys   any
	mode  fs.FileMode
	isDir bool

	ModeString    string `json:"mode"`
	ModeOctString string `json:"mode_oct"`

	F_name string `json:"name"`

	I_size int64 `json:"size"`

	basename    string
	ext         string
	suffix      string
	dir         string
	curDir      string
	isAbs       bool
	dirAbs      string
	pathAbs     string
	pathDisplay string

	// S_time              string `json:"time"`
	S_mtime string `json:"mtime"`
	S_ctime string `json:"ctime"`
	S_atime string `json:"atime"`

	// _time *time.Time
	ctime *time.Time
	atime *time.Time
	mtime *time.Time

	F_owner     string `json:"owner"`
	F_group     string `json:"group"`
	F_inode     uint64 `json:"inode"`
	F_hardLinks uint64 `json:"hard_links"`
	F_blocks    int64  `json:"blocks"`

	F_deviceNumbers string // `json:""`
}

func (fi *FakeFileInfo) Prepare() error {
	if fi.ModeString != "" {
		mode, err := parse.ParseMode(fi.ModeString)
		if err != nil {
			return err
		}
		fi.mode = mode
	} else if fi.ModeOctString != "" {
		modeOct, err := strconv.ParseInt(fi.ModeOctString, 8, 16)
		if err != nil {
			return err
		}
		fi.mode |= fs.FileMode(modeOct)
	}
	mode := fi.mode
	if mode&fs.ModeDir > 0 {
		fi.isDir = true
	}
	name := fi.F_name
	if strings.HasSuffix(name, "/") {
		name = name[:len(name)-1]
		fi.isDir = true
		fi.mode |= fs.ModeDir
	}
	pname := fspaths.SplitExt(name)
	fi.basename = pname.Base
	fi.ext = pname.Ext
	fi.suffix = pname.Suffix
	fi.dir = fspaths.Dir(name)
	// fi.curDir
	fi.isAbs = fspaths.IsAbs(name)
	dirAbs, err := fspaths.Abs(fi.dir)
	if err != nil {
		return err
	}
	fi.dirAbs = dirAbs
	pathAbs, err := fspaths.Abs(name)
	if err != nil {
		return err
	}
	fi.pathAbs = pathAbs
	fi.pathDisplay = name
	if fi.S_mtime != "" {
		_time, err := time.Parse(timeFmt, fi.S_mtime)
		if err != nil {
			return err
		}
		fi.mtime = &_time
	}
	if fi.S_ctime != "" {
		_time, err := time.Parse(timeFmt, fi.S_ctime)
		if err != nil {
			return err
		}
		fi.ctime = &_time
	}
	if fi.S_atime != "" {
		_time, err := time.Parse(timeFmt, fi.S_atime)
		if err != nil {
			return err
		}
		fi.atime = &_time
	}
	return nil
}

func (fi *FakeFileInfo) Name() string {
	return fi.F_name
}

func (fi *FakeFileInfo) Size() int64 {
	return fi.I_size
}

func (fi *FakeFileInfo) Mode() fs.FileMode {
	return fi.mode
}

func (fi *FakeFileInfo) ModTime() time.Time {
	if fi.mtime != nil {
		return *fi.mtime
	}
	return time.Time{}
}

func (fi *FakeFileInfo) IsDir() bool {
	return fi.isDir
}

func (fi *FakeFileInfo) Sys() any {
	return fi.sys
}

func (fi *FakeFileInfo) Basename() string {
	return fi.basename
}

func (fi *FakeFileInfo) Ext() string {
	return fi.ext
}

func (fi *FakeFileInfo) Suffix() string {
	return fi.suffix
}

func (fi *FakeFileInfo) Dir() string {
	return fi.dir
}

func (fi *FakeFileInfo) CurDir() string {
	return fi.curDir
}

func (fi *FakeFileInfo) IsAbs() bool {
	return fi.isAbs
}

func (fi *FakeFileInfo) DirAbs() string {
	return fi.dirAbs
}

func (fi *FakeFileInfo) PathAbs() string {
	return fi.pathAbs
}

func (fi *FakeFileInfo) PathDisplay() string {
	return fi.pathDisplay
}

func (fi *FakeFileInfo) Time(colName string) *time.Time {
	switch colName {
	case "mtime":
		return fi.mtime
	case "ctime":
		return fi.ctime
	case "atime":
		return fi.atime
	}
	return nil
}

func (fi *FakeFileInfo) Owner() string {
	return fi.F_owner
}

func (fi *FakeFileInfo) Group() string {
	return fi.F_group
}

func (fi *FakeFileInfo) Inode() (uint64, error) {
	return fi.F_inode, nil
}

func (fi *FakeFileInfo) NumberOfHardLinks() (uint64, error) {
	return fi.F_hardLinks, nil
}

func (fi *FakeFileInfo) DeviceNumbers() (string, error) {
	return fi.F_deviceNumbers, nil
}

func (fi *FakeFileInfo) CTime() *time.Time {
	return fi.ctime
}

func (fi *FakeFileInfo) ATime() *time.Time {
	return fi.atime
}

func (fi *FakeFileInfo) Blocks() int64 {
	return fi.F_blocks
}
