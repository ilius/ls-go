//go:build darwin

package lsplatform

import (
	"strconv"
	"syscall"

	"github.com/ilius/ls-go/lsplatform/unix"
)

func (*LocalPlatform) DeviceNumbers(info FileInfo) (string, error) {
	stat := info.Sys().(*syscall.Stat_t)
	major := strconv.FormatInt(int64(unix.Major(uint64(stat.Rdev))), 10)
	minor := strconv.FormatInt(int64(unix.Minor(uint64(stat.Rdev))), 10)
	return major + "," + minor, nil
}
