//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || zos

package lsplatform

import (
	"strconv"
	"syscall"

	"github.com/ilius/ls-go/lsplatform/unix"
)

func (*LocalPlatform) DeviceNumbers(info FileInfo) (string, error) {
	stat := info.Sys().(*syscall.Stat_t)
	major := strconv.FormatInt(int64(unix.Major(stat.Rdev)), 10)
	minor := strconv.FormatInt(int64(unix.Minor(stat.Rdev)), 10)
	return major + "," + minor, nil
}
