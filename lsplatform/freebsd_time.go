//go:build freebsd || darwin || netbsd

package lsplatform

import (
	"syscall"
	"time"
)

func (*LocalPlatform) FileCTime(fileInfo FileInfo) *time.Time {
	stat := fileInfo.Sys().(*syscall.Stat_t)
	ctime := time.Unix(int64(stat.Ctimespec.Sec), int64(stat.Ctimespec.Nsec))
	return &ctime
}

func (*LocalPlatform) FileATime(fileInfo FileInfo) *time.Time {
	stat := fileInfo.Sys().(*syscall.Stat_t)
	atime := time.Unix(int64(stat.Atimespec.Sec), int64(stat.Atimespec.Nsec))
	return &atime
}
