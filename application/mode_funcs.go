package application

import (
	"io/fs"
	"os"
	"strings"

	"github.com/ilius/ls-go/lscolors"
)

func fileTypeSymbol(mode fs.FileMode) string {
	// this "type" is not the file extension, but type as far as the OS is concerned
	if mode&os.ModeDir != 0 {
		return "d"
	}
	if mode&os.ModeSymlink != 0 {
		return "l"
	}
	if mode&os.ModeDevice != 0 {
		if mode&os.ModeCharDevice == 0 {
			return "b" // block device
		}
		return "c" // character device
	}
	if mode&os.ModeNamedPipe != 0 {
		return "p"
	}
	if mode&os.ModeSocket != 0 {
		return "s"
	}
	return "-"
}

func specialPermSymbol(bits fs.FileMode, i uint) string {
	if i == 0 {
		if bits&1 == 0 {
			return "T" // sticky but non-executable
		}
		return "t" // sticky
	}
	if bits&1 == 0 {
		return "S" // SUID/SGID but non-executable
	}
	return "s" // SUID/SGID
}

func rwxString(mode fs.FileMode, i uint, color *lscolors.Style) string {
	bits := mode >> (i * 3)
	parts := []string{}
	if bits&4 == 0 {
		parts = append(parts, "-")
	} else {
		parts = append(parts, "r")
	}
	if bits&2 == 0 {
		parts = append(parts, "-")
	} else {
		parts = append(parts, "w")
	}
	if mode&specialPermModes[i] != 0 {
		// SUID/SGID/sticky
		parts = append(parts, specialPermSymbol(bits, i))
	} else {
		if bits&1 == 0 {
			parts = append(parts, "-")
		} else {
			parts = append(parts, "x")
		}
	}
	str := strings.Join(parts, "")
	if color != nil {
		str = app.Colorize(str, color)
	}
	return str
}

func formatModeNoColor(info fs.FileInfo) string {
	// info.Mode().String() does not produce the same output as `ls`, so we must build that string manually
	mode := info.Mode()
	return strings.Join([]string{
		fileTypeSymbol(mode),
		rwxString(mode, 2, nil),
		rwxString(mode, 1, nil),
		rwxString(mode, 0, nil),
	}, "")
}
