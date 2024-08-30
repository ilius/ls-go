//go:build !windows

package lsplatform

import (
	"io"
	"os"
	"os/user"
	"strconv"
	"syscall"

	etcgroup "github.com/ilius/etc/group"
	"github.com/ilius/etc/passwd"
)

var (
	userNameById  = map[uint32]string{}
	groupNameById = map[uint32]string{}

	systemUserNames = []string{}
)

func init() {
	userMap, err := passwd.Parse()
	if err != nil {
		panic(err)
	}
	for username, user := range userMap {
		uid, err := strconv.ParseInt(user.UID, 10, 32)
		if err != nil {
			panic(err)
		}
		userNameById[uint32(uid)] = username
		if uid > 0 && uid < 1000 {
			systemUserNames = append(systemUserNames, username)
		}
	}
	groupMap, err := etcgroup.Parse()
	if err != nil {
		panic(err)
	}
	for groupName, group := range groupMap {
		gid, err := strconv.ParseInt(group.GID, 10, 32)
		if err != nil {
			panic(err)
		}
		groupNameById[uint32(gid)] = groupName
	}
}

func lookupUserId(uid uint32) string {
	name, ok := userNameById[uid]
	if ok {
		return name
	}
	return strconv.FormatUint(uint64(uid), 10)
}

func lookupGroupId(gid uint32) string {
	name, ok := groupNameById[gid]
	if ok {
		return name
	}
	gids := strconv.FormatUint(uint64(gid), 10)
	g, err := user.LookupGroupId(gids)
	if err != nil {
		groupNameById[gid] = gids
		return gids
	}
	groupNameById[gid] = g.Name
	return g.Name
}

// RootUserName returns name of root user (the main admin)
func (*LocalPlatform) RootUserName() string {
	return "root"
}

func (*LocalPlatform) SystemUserNames() []string {
	return systemUserNames
}

// UserName returns name of current user
func (*LocalPlatform) UserName() string {
	return os.Getenv("USER")
}

func (*LocalPlatform) OwnerAndGroupNames(fileInfo FileInfo) (*OwnerGroup, error) {
	stat := fileInfo.Sys().(*syscall.Stat_t)
	return &OwnerGroup{
		lookupUserId(stat.Uid),
		lookupGroupId(stat.Gid),
	}, nil
}

func (*LocalPlatform) OwnerAndGroupIDs(fileInfo FileInfo) (*OwnerGroup, error) {
	stat := fileInfo.Sys().(*syscall.Stat_t)
	return &OwnerGroup{
		strconv.FormatUint(uint64(stat.Uid), 10),
		strconv.FormatUint(uint64(stat.Gid), 10),
	}, nil
}

func (*LocalPlatform) NumberOfHardLinks(fileInfo FileInfo) (uint64, error) {
	if sys := fileInfo.Sys(); sys != nil {
		if stat, ok := sys.(*syscall.Stat_t); ok {
			return uint64(stat.Nlink), nil
		}
	}
	return 0, nil
}

func (*LocalPlatform) FileInode(fileInfo FileInfo) (uint64, error) {
	return fileInfo.Sys().(*syscall.Stat_t).Ino, nil
}

// FileBlocks returns number of 1024-byte blocks occupied by a file
func (*LocalPlatform) FileBlocks(fileInfo FileInfo) int64 {
	return fileInfo.Sys().(*syscall.Stat_t).Blocks / 2
}

func (*LocalPlatform) EmptyFileInfoSys() any {
	return &syscall.Stat_t{}
}

func (*LocalPlatform) OutputAndError(colors bool) (io.Writer, io.Writer) {
	return os.Stdout, os.Stderr
}
