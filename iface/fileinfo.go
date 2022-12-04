package iface

import (
	"io/fs"
	"time"
)

type FileInfo interface {
	fs.FileInfo

	Basename() string
	Ext() string
	Suffix() string
	Dir() string
	CurDir() string
	IsAbs() bool
	DirAbs() string
	PathAbs() string
	PathDisplay() string
	Time(colName string) *time.Time
	Owner() string
	Group() string
	Inode() (uint64, error)
	NumberOfHardLinks() (uint64, error)
	DeviceNumbers() (string, error)
	CTime() *time.Time
	ATime() *time.Time
	Blocks() int64
}
