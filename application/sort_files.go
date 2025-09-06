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
		sort.Sort(sort.Reverse(SizeSorter(files)))
	} else {
		sort.Sort(SizeSorter(files))
	}
}

func sortByFileSize(files []*DisplayItem, reverse bool) {
	if reverse {
		sort.Sort(sort.Reverse(FileSizeSorter(files)))
	} else {
		sort.Sort(FileSizeSorter(files))
	}
}

func sortByTime(files []*DisplayItem, reverse bool) {
	if reverse {
		sort.Sort(sort.Reverse(TimeSorter(files)))
	} else {
		sort.Sort(TimeSorter(files))
	}
}

func sortByExtension(files []*DisplayItem, reverse bool) {
	if reverse {
		sort.Sort(sort.Reverse(ExtensionSorter(files)))
	} else {
		sort.Sort(ExtensionSorter(files))
	}
}

func sortByKind(files []*DisplayItem, reverse bool) {
	if reverse {
		sort.Sort(sort.Reverse(KindSorter(files)))
	} else {
		sort.Sort(KindSorter(files))
	}
}

func sortByInode(files []*DisplayItem, reverse bool) {
	if reverse {
		sort.Sort(sort.Reverse(InodeSorter(files)))
	} else {
		sort.Sort(InodeSorter(files))
	}
}

func sortByLinks(files []*DisplayItem, reverse bool) {
	if reverse {
		sort.Sort(sort.Reverse(HardLinksSorter(files)))
	} else {
		sort.Sort(HardLinksSorter(files))
	}
}

func sortByMode(files []*DisplayItem, reverse bool) {
	if reverse {
		sort.Sort(sort.Reverse(ModeSorter(files)))
	} else {
		sort.Sort(ModeSorter(files))
	}
}

func sortByNameLen(files []*DisplayItem, reverse bool) {
	if reverse {
		sort.Sort(sort.Reverse(NameLengthSorter(files)))
	} else {
		sort.Sort(NameLengthSorter(files))
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
			sort.Sort(sort.Reverse(DirContentsCountSorter(dirs)))
			return
		}
		sort.Sort(DirContentsCountSorter(dirs))
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
