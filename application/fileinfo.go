package application

import (
	"fmt"
	"io/fs"
	"time"

	c "github.com/ilius/ls-go/common"
	"github.com/ilius/ls-go/escape"
	"github.com/ilius/ls-go/iface"
)

type FileInfo = iface.FileInfo

type FileInfoImp struct {
	fs.FileInfo

	owner *string
	group *string

	basename string
	ext      string
	suffix   string

	dir    string
	curDir string
	isAbs  bool
}

func (info *FileInfoImp) Basename() string {
	return info.basename
}

func (info *FileInfoImp) Ext() string {
	return info.ext
}

func (info *FileInfoImp) Suffix() string {
	return info.suffix
}

func (info *FileInfoImp) Dir() string {
	return info.dir
}

func (info *FileInfoImp) CurDir() string {
	return info.curDir
}

func (info *FileInfoImp) IsAbs() bool {
	return info.isAbs
}

func (info *FileInfoImp) DirAbs() string {
	dirAbs, err := app.FileSystem.Abs(info.dir)
	check(err)
	return dirAbs
}

func (info *FileInfoImp) PathAbs() string {
	return app.FileSystem.Join(info.DirAbs(), info.Name())
}

func (info *FileInfoImp) quoteAndEscape(path string) string {
	path = quoteFileName(path)
	if app.EnsureASCII {
		path = escape.EscapeToASCII(path)
	}
	return path
}

func (info *FileInfoImp) PathDisplay() string {
	name := info.Name()
	curDir := info.curDir
	if curDir == "" {
		curDir = app.workDir
	}
	dirAbs := info.DirAbs()
	if dirAbs == curDir {
		return info.quoteAndEscape(name)
	}
	name = app.FileSystem.Join(dirAbs, name)
	if info.isAbs {
		return info.quoteAndEscape(name)
	}
	name, err := app.FileSystem.Rel(curDir, name)
	check(err)
	return info.quoteAndEscape(name)
}

func (info *FileInfoImp) Time(colName string) *time.Time {
	switch colName {
	case c.C_MTime:
		_time := info.ModTime()
		return &_time
	case c.C_CTime:
		return app.Platform.FileCTime(info)
	case c.C_ATime:
		return app.Platform.FileATime(info)
	}
	panic(fmt.Errorf("invalid colName=%#v", colName))
}

const unknown = "unknown"

func (info *FileInfoImp) Owner() string {
	if info.owner != nil {
		return *info.owner
	}
	og, err := getOwnerAndGroup(info)
	if err != nil {
		app.AddError(err)
		str := unknown
		info.owner = &str
		info.group = &str
		return unknown
	}
	info.owner = &og.Owner
	info.group = &og.Group
	return og.Owner
}

func (info *FileInfoImp) Group() string {
	if info.group != nil {
		return *info.group
	}
	og, err := getOwnerAndGroup(info)
	if err != nil {
		app.AddError(err)
		str := unknown
		info.owner = &str
		info.group = &str
		return unknown
	}
	return og.Group
}

func (info *FileInfoImp) Inode() (uint64, error) {
	return app.Platform.FileInode(info)
}

func (info *FileInfoImp) NumberOfHardLinks() (uint64, error) {
	return app.Platform.NumberOfHardLinks(info)
}

func (info *FileInfoImp) DeviceNumbers() (string, error) {
	return app.Platform.DeviceNumbers(info)
}

func (info *FileInfoImp) CTime() *time.Time {
	return app.Platform.FileCTime(info)
}

func (info *FileInfoImp) ATime() *time.Time {
	return app.Platform.FileATime(info)
}

func (info *FileInfoImp) Blocks() int64 {
	return app.Platform.FileBlocks(info)
}

type FileInfoLow struct {
	modTime time.Time
	sys     any
	name    string
	size    int64
	mode    fs.FileMode
	isDir   bool
}

func (fi *FileInfoLow) Name() string {
	return fi.name
}

func (fi *FileInfoLow) Size() int64 {
	return fi.size
}

func (fi *FileInfoLow) Mode() fs.FileMode {
	return fi.mode
}

func (fi *FileInfoLow) ModTime() time.Time {
	return fi.modTime
}

func (fi *FileInfoLow) IsDir() bool {
	return fi.isDir
}

func (fi *FileInfoLow) Sys() any {
	return fi.sys
}

func NewFileInfoLowFrom(fi fs.FileInfo) *FileInfoLow {
	return &FileInfoLow{
		name:    fi.Name(),
		size:    fi.Size(),
		mode:    fi.Mode(),
		modTime: fi.ModTime(),
		isDir:   fi.IsDir(),
		sys:     fi.Sys(),
	}
}
