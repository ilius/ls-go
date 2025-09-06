package application

import (
	"sort"

	c "github.com/ilius/ls-go/common"
)

func sortFiles(files []*DisplayItem, col string, reverse bool) {
	switch col {
	case c.S_NONE:
		break // do not sort
	case c.S_NAME:
		sortByName(files, reverse)
	case c.S_BASENAME:
		sortByBasename(files, reverse)
	case c.S_SIZE:
		sortBySize(files, reverse)
	case c.S_FILESIZE:
		sortByFileSize(files, reverse)
	case c.S_TIME:
		sortByTime(files, reverse)
	// case S_VERSION:
	case c.S_EXTENSION:
		sortByExtension(files, reverse)
	case c.S_KIND:
		sortByKind(files, reverse)
	case c.S_INODE:
		sortByInode(files, reverse)
	case c.S_LINKS:
		sortByLinks(files, reverse)
	case c.S_MODE:
		sortByMode(files, reverse)
	case c.S_NAME_LEN:
		sortByNameLen(files, reverse)
	default: // default is (basename, extension)
		sortDefault(files, reverse)
	}
}

func sortByName(files []*DisplayItem, reverse bool) {
	if reverse {
		sort.Sort(sort.Reverse(NameSorter(files)))
	} else {
		sort.Sort(NameSorter(files))
	}
}

func sortByBasename(files []*DisplayItem, reverse bool) {
	if reverse {
		sort.Sort(sort.Reverse(BasenameSorter(files)))
	} else {
		sort.Sort(BasenameSorter(files))
	}
}

func sortBySize(files []*DisplayItem, reverse bool) {
	if reverse {
		sort.Sort(sort.Reverse(BySize(files)))
	} else {
		sort.Sort(BySize(files))
	}
}

func sortByFileSize(files []*DisplayItem, reverse bool) {
	if reverse {
		sort.Sort(sort.Reverse(ByFileSize(files)))
	} else {
		sort.Sort(ByFileSize(files))
	}
}

func sortByTime(files []*DisplayItem, reverse bool) {
	if reverse {
		sort.Sort(sort.Reverse(ByTime(files)))
	} else {
		sort.Sort(ByTime(files))
	}
}

func sortByExtension(files []*DisplayItem, reverse bool) {
	if reverse {
		sort.Sort(sort.Reverse(ByExtension(files)))
	} else {
		sort.Sort(ByExtension(files))
	}
}

func sortByKind(files []*DisplayItem, reverse bool) {
	if reverse {
		sort.Sort(sort.Reverse(ByKind(files)))
	} else {
		sort.Sort(ByKind(files))
	}
}

func sortByInode(files []*DisplayItem, reverse bool) {
	if reverse {
		sort.Sort(sort.Reverse(ByInode(files)))
	} else {
		sort.Sort(ByInode(files))
	}
}

func sortByLinks(files []*DisplayItem, reverse bool) {
	if reverse {
		sort.Sort(sort.Reverse(ByHardLinks(files)))
	} else {
		sort.Sort(ByHardLinks(files))
	}
}

func sortByMode(files []*DisplayItem, reverse bool) {
	if reverse {
		sort.Sort(sort.Reverse(ByMode(files)))
	} else {
		sort.Sort(ByMode(files))
	}
}

func sortByNameLen(files []*DisplayItem, reverse bool) {
	if reverse {
		sort.Sort(sort.Reverse(ByNameLength(files)))
	} else {
		sort.Sort(ByNameLength(files))
	}
}

// default is (basename, extension)
func sortDefault(files []*DisplayItem, reverse bool) {
	if reverse {
		sort.Sort(sort.Reverse(DefaultSorter(files)))
	} else {
		sort.Sort(DefaultSorter(files))
	}
}

func sortDirs(dirs []*DisplayItem, col string, reverse bool) {
	if len(dirs) == 0 {
		return
	}
	switch col {
	case c.S_SIZE:
		if reverse {
			sort.Sort(sort.Reverse(ByDirContentsCount(dirs)))
			return
		}
		sort.Sort(ByDirContentsCount(dirs))
		return
	case c.S_FILESIZE:
		if reverse {
			sort.Sort(sort.Reverse(DefaultSorter(dirs)))
			return
		}
		sort.Sort(DefaultSorter(dirs))
		return
	}
	sortFiles(dirs, col, reverse)
}
