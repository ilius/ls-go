package application

import (
	"sort"

	. "github.com/ilius/ls-go/common"
)

func sortFiles(files []*DisplayItem, col string, reverse bool) {
	switch col {
	case S_NONE:
		break // do not sort
	case S_SIZE:
		if reverse {
			sort.Sort(sort.Reverse(BySize(files)))
		} else {
			sort.Sort(BySize(files))
		}
	case S_FILESIZE:
		if reverse {
			sort.Sort(sort.Reverse(ByFileSize(files)))
		} else {
			sort.Sort(ByFileSize(files))
		}
	case S_TIME:
		if reverse {
			sort.Sort(sort.Reverse(ByTime(files)))
		} else {
			sort.Sort(ByTime(files))
		}
	// case S_VERSION:
	case S_EXTENSION:
		if reverse {
			sort.Sort(sort.Reverse(ByExtension(files)))
		} else {
			sort.Sort(ByExtension(files))
		}
	case S_KIND:
		if reverse {
			sort.Sort(sort.Reverse(ByKind(files)))
		} else {
			sort.Sort(ByKind(files))
		}
	case S_INODE:
		if reverse {
			sort.Sort(sort.Reverse(ByInode(files)))
		} else {
			sort.Sort(ByInode(files))
		}
	case S_LINKS:
		if reverse {
			sort.Sort(sort.Reverse(ByHardLinks(files)))
		} else {
			sort.Sort(ByHardLinks(files))
		}
	case S_MODE:
		if reverse {
			sort.Sort(sort.Reverse(ByMode(files)))
		} else {
			sort.Sort(ByMode(files))
		}
	case S_NAME_LEN:
		if reverse {
			sort.Sort(sort.Reverse(ByNameLength(files)))
		} else {
			sort.Sort(ByNameLength(files))
		}
	default:
		if reverse {
			sort.Sort(sort.Reverse(DefaultSorter(files)))
		} else {
			sort.Sort(DefaultSorter(files))
		}
	}
}

func sortDirs(dirs []*DisplayItem, col string, reverse bool) {
	if len(dirs) == 0 {
		return
	}
	switch col {
	case S_SIZE:
		if reverse {
			sort.Sort(sort.Reverse(ByDirContentsCount(dirs)))
			return
		}
		sort.Sort(ByDirContentsCount(dirs))
		return
	case S_FILESIZE:
		if reverse {
			sort.Sort(sort.Reverse(DefaultSorter(dirs)))
			return
		}
		sort.Sort(DefaultSorter(dirs))
		return
	}
	sortFiles(dirs, col, reverse)
}
