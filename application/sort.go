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

type NumericNameSorter ItemSorter

func (s NumericNameSorter) Len() int      { return len(s) }
func (s NumericNameSorter) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s NumericNameSorter) Less(i, j int) bool {
	i_base := strings.ToLower(s[i].Basename())
	j_base := strings.ToLower(s[j].Basename())

	if i_base != j_base {
		return i_base < j_base
	}
	return s[i].Ext() < s[j].Ext()
}

// sort by size in decending order
type BySize ItemSorter

func (s BySize) Len() int      { return len(s) }
func (s BySize) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s BySize) Less(i, j int) bool {
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

// ByFileSize sorts files by size, and directories by lowercased name
// and sort directories after files
type ByFileSize ItemSorter

func (s ByFileSize) Len() int      { return len(s) }
func (s ByFileSize) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s ByFileSize) Less(i, j int) bool {
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
type ByDirContentsCount ItemSorter

func (s ByDirContentsCount) Len() int      { return len(s) }
func (s ByDirContentsCount) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s ByDirContentsCount) Less(i, j int) bool {
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
type ByTime ItemSorter

func (s ByTime) Len() int      { return len(s) }
func (s ByTime) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s ByTime) Less(i, j int) bool {
	// do NOT compare unix times (returned by _time.Unix)
	// because of DST stuff, it's complicated
	tm1 := s[i].Time
	tm2 := s[j].Time
	if tm1 == nil || tm2 == nil {
		panic("time is nil")
		// return false
	}
	return (*tm1).After(*tm2)
}

type ByExtension ItemSorter

func (s ByExtension) Len() int      { return len(s) }
func (s ByExtension) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s ByExtension) Less(i, j int) bool {
	return s[i].Ext() < s[j].Ext()
}

// ByKind tells `sort.Sort` how to sort by file extension
type ByKind ItemSorter

func (s ByKind) Len() int      { return len(s) }
func (s ByKind) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s ByKind) Less(i, j int) bool {
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

type ByInode ItemSorter

func (s ByInode) Len() int      { return len(s) }
func (s ByInode) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s ByInode) Less(i, j int) bool {
	inode1, _ := app.Platform.FileInode(s[i])
	inode2, _ := app.Platform.FileInode(s[j])
	return inode1 < inode2
}

type ByHardLinks ItemSorter

func (s ByHardLinks) Len() int      { return len(s) }
func (s ByHardLinks) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s ByHardLinks) Less(i, j int) bool {
	n1, _ := app.Platform.NumberOfHardLinks(s[i])
	n2, _ := app.Platform.NumberOfHardLinks(s[j])
	return n1 > n2
}

type ByMode ItemSorter

func (s ByMode) Len() int      { return len(s) }
func (s ByMode) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s ByMode) Less(i, j int) bool {
	return s[i].Mode() > s[j].Mode()
}

type ByNameLength ItemSorter

func (s ByNameLength) Len() int      { return len(s) }
func (s ByNameLength) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s ByNameLength) Less(i, j int) bool {
	return len(s[i].Name()) < len(s[j].Name())
}
