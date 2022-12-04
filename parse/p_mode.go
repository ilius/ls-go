package parse

import (
	"fmt"
	"io/fs"
)

var specialPermModes = [3]fs.FileMode{
	fs.ModeSticky,
	fs.ModeSetgid,
	fs.ModeSetuid,
}

func parsePerms(str string, index int) (fs.FileMode, error) {
	perm := fs.FileMode(0)
	switch str[0] {
	case 'r':
		perm = 4
	case '-':
	default:
		return 0, fmt.Errorf("invalid perm=%#v", str)
	}
	switch str[1] {
	case 'w':
		perm += 2
	case '-':
	default:
		return 0, fmt.Errorf("invalid perm=%#v", str)
	}
	switch str[2] {
	case 'x':
		perm += 1
	case '-':
	case 's':
		if index == 2 {
			return 0, fmt.Errorf("invalid perm=%#v", str)
		}
		perm += 1 + specialPermModes[index]
	case 'S':
		if index == 2 {
			return 0, fmt.Errorf("invalid perm=%#v", str)
		}
		perm += specialPermModes[index]
	case 't':
		if index != 2 {
			return 0, fmt.Errorf("invalid perm=%#v", str)
		}
		perm += 1 + specialPermModes[index]
	case 'T':
		if index != 2 {
			return 0, fmt.Errorf("invalid perm=%#v", str)
		}
		perm += specialPermModes[index]
	default:
		return 0, fmt.Errorf("invalid perm=%#v", str)
	}
	return perm, nil
}

func ParseMode(str string) (fs.FileMode, error) {
	if len(str) != 10 {
		return 0, fmt.Errorf("invalid mode=%#v", str)
	}
	mode := fs.FileMode(0)
	switch str[0] {
	case '-':
	case 'd':
		mode = fs.ModeDir
	case 'l':
		mode = fs.ModeSymlink
	case 'b':
		mode = fs.ModeDevice
	case 'c':
		mode = fs.ModeDevice | fs.ModeCharDevice
	case 'p':
		mode = fs.ModeNamedPipe
	case 's':
		mode = fs.ModeSocket
	}
	user, err := parsePerms(str[1:4], 0)
	if err != nil {
		return 0, nil
	}
	group, err := parsePerms(str[4:7], 1)
	if err != nil {
		return 0, nil
	}
	others, err := parsePerms(str[7:], 2)
	if err != nil {
		return 0, nil
	}
	mode |= user<<6 | group<<3 | others
	return mode, nil
}
