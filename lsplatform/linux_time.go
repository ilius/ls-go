//go:build linux || openbsd || dragonfly || solaris || android
// +build linux openbsd dragonfly solaris android

package lsplatform

import (
	"syscall"
	"time"
)

func (*LocalPlatform) FileCTime(fileInfo FileInfo) *time.Time {
	stat := fileInfo.Sys().(*syscall.Stat_t)
	ctime := time.Unix(int64(stat.Ctim.Sec), int64(stat.Ctim.Nsec))
	return &ctime
}

func (*LocalPlatform) FileATime(fileInfo FileInfo) *time.Time {
	stat := fileInfo.Sys().(*syscall.Stat_t)
	atime := time.Unix(int64(stat.Atim.Sec), int64(stat.Atim.Nsec))
	return &atime
}
