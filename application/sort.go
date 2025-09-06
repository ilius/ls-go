package application

import (
	"strings"
)

type ItemSorter []*DisplayItem

// DefaultSorter is the default sorter for files and directories
// it first sorts by lowercased basename, then by extension
type DefaultSorter ItemSorter

func (s DefaultSorter) Len() int      { return len(s) }
func (s DefaultSorter) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s DefaultSorter) Less(i, j int) bool {
	i_base := strings.ToLower(s[i].Basename())
	j_base := strings.ToLower(s[j].Basename())
	if i_base != j_base {
		return i_base < j_base
	}
	return s[i].Ext() < s[j].Ext()
}

// NameSorter sorts by full name
type NameSorter ItemSorter

func (s NameSorter) Len() int      { return len(s) }
func (s NameSorter) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s NameSorter) Less(i, j int) bool {
	return strings.ToLower(s[i].Name()) < strings.ToLower(s[j].Name())
}

type BasenameSorter ItemSorter

func (s BasenameSorter) Len() int      { return len(s) }
func (s BasenameSorter) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s BasenameSorter) Less(i, j int) bool {
	return strings.ToLower(s[i].Basename()) < strings.ToLower(s[j].Basename())
}

// sort by size in decending order
type SizeSorter ItemSorter

func (s SizeSorter) Len() int      { return len(s) }
func (s SizeSorter) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s SizeSorter) Less(i, j int) bool {
	info1 := s[i]
	info2 := s[j]
	if info1.IsDir() && info2.IsDir() {
		n1, err := app.FileSystem.CountDirContents(info1.PathAbs())
		check(err)
		n2, _ := app.FileSystem.CountDirContents(info2.PathAbs())
		check(err)
		return n1 > n2
	}
	return info1.Size() > info2.Size()
}

// FileSizeSorter sorts files by size, and directories by lowercased name
// and sort directories after files
type FileSizeSorter ItemSorter

func (s FileSizeSorter) Len() int      { return len(s) }
func (s FileSizeSorter) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s FileSizeSorter) Less(i, j int) bool {
	info1 := s[i]
	info2 := s[j]
	if info1.IsDir() {
		if !info2.IsDir() {
			return false
		}
		return strings.ToLower(info1.Name()) < strings.ToLower(info2.Name())
	}
	if info2.IsDir() {
		return true
	}
	return info1.Size() > info2.Size()
}

// sorts by the number of contents (files and directories) in directories
// directories always come before files
// and remember, this includes hidden contents as well, so to count
// manually, use `ls -A1` or `ls-go -a1`
type DirContentsCountSorter ItemSorter

func (s DirContentsCountSorter) Len() int      { return len(s) }
func (s DirContentsCountSorter) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s DirContentsCountSorter) Less(i, j int) bool {
	info1 := s[i]
	info2 := s[j]
	if !info1.IsDir() {
		return false
	}
	if !info2.IsDir() {
		return true
	}
	n1, _ := app.FileSystem.CountDirContents(info1.PathAbs())
	n2, _ := app.FileSystem.CountDirContents(info2.PathAbs())
	return n1 > n2
}

// sort by time (modified time by default) in decending order (newer first)
type TimeSorter ItemSorter

func (s TimeSorter) Len() int      { return len(s) }
func (s TimeSorter) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s TimeSorter) Less(i, j int) bool {
	// do NOT compare unix times (returned by _time.Unix)
	// because of DST stuff, it's complicated
	tm1 := s[i].Time
	tm2 := s[j].Time
	if tm1 == nil || tm2 == nil {
		panic("time is nil")
		// return false
	}
	return tm1.After(*tm2)
}

type ExtensionSorter ItemSorter

func (s ExtensionSorter) Len() int      { return len(s) }
func (s ExtensionSorter) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s ExtensionSorter) Less(i, j int) bool {
	return s[i].Ext() < s[j].Ext()
}

// KindSorter tells `sort.Sort` how to sort by file extension
type KindSorter ItemSorter

func (s KindSorter) Len() int      { return len(s) }
func (s KindSorter) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s KindSorter) Less(i, j int) bool {
	var kindi, kindj string
	if s[i].Basename() == "" {
		kindi = "."
	} else if s[i].Ext() == "" {
		kindi = "0"
	} else {
		kindi = s[i].Ext()
	}
	if s[j].Basename() == "" {
		kindj = "."
	} else if s[j].Ext() == "" {
		kindj = "0"
	} else {
		kindj = s[j].Ext()
	}
	if kindi == kindj {
		if kindi == "." {
			return s[i].Ext() < s[j].Ext()
		}
		return s[i].Basename() < s[j].Basename()
	}
	return kindi < kindj
}

type InodeSorter ItemSorter

func (s InodeSorter) Len() int      { return len(s) }
func (s InodeSorter) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s InodeSorter) Less(i, j int) bool {
	inode1, _ := app.Platform.FileInode(s[i])
	inode2, _ := app.Platform.FileInode(s[j])
	return inode1 < inode2
}

type HardLinksSorter ItemSorter

func (s HardLinksSorter) Len() int      { return len(s) }
func (s HardLinksSorter) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s HardLinksSorter) Less(i, j int) bool {
	n1, _ := app.Platform.NumberOfHardLinks(s[i])
	n2, _ := app.Platform.NumberOfHardLinks(s[j])
	return n1 > n2
}

type ModeSorter ItemSorter

func (s ModeSorter) Len() int      { return len(s) }
func (s ModeSorter) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s ModeSorter) Less(i, j int) bool {
	return s[i].Mode() > s[j].Mode()
}

type NameLengthSorter ItemSorter

func (s NameLengthSorter) Len() int      { return len(s) }
func (s NameLengthSorter) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s NameLengthSorter) Less(i, j int) bool {
	return len(s[i].Name()) < len(s[j].Name())
}
