package application

import (
	"fmt"
	"os"

	. "github.com/ilius/ls-go/common"
)

type FileNameGetterPlain struct {
	*FileNameParams
}

func (f *FileNameGetterPlain) Value(item any) (any, error) {
	info, ok := item.(FileInfo)
	if !ok {
		return "", fmt.Errorf("Value: invalid type %T, must be FileInfo", item)
	}
	parentDirAbs := info.DirAbs()
	// read some info about linked file if this item is a symlink
	var link *LinkInfo
	if info.Mode()&os.ModeSymlink != 0 {
		link = getLinkInfo(info, parentDirAbs, f.linkRel)
	}
	displayName := f.nameString(info, link)

	if f.showLinks && info.Mode()&os.ModeSymlink != 0 {
		displayName += " ► " + f.linkTargetString(link)
	}

	return displayName, nil
}

func (f *FileNameGetterPlain) ValueString(colName string, item any) (string, error) {
	info, ok := item.(FileInfo)
	if !ok {
		return "", fmt.Errorf("ValueString: invalid type %T, must be FileInfo", item)
	}
	filename := info.Name()
	if f.fullPath || info.IsAbs() {
		// info.Dir for relative path, info.DirAbs() for absoulte path
		filename = app.FileSystem.Join(info.DirAbs(), filename)
	}
	if info.IsDir() {
		filename += "/"
	}
	str, err := app.FormatValue(colName, filename)
	if err != nil {
		return "", err
	}

	if f.showLinks {
		target := ""
		if info.Mode()&os.ModeSymlink != 0 {
			link := getLinkInfo(info, info.DirAbs(), f.linkRel)
			if link != nil {
				target = link.target
			}
		}
		targetFormatted, err := app.FormatValue(C_LinkTarget, target)
		if err != nil {
			return "", err
		}
		str += f.showLinksSep + targetFormatted
	}
	return str, nil
}

func (f *FileNameGetterPlain) Format(item any, value any) (string, error) {
	return value.(string), nil
}

func (f *FileNameGetterPlain) nameString(info FileInfo, link *LinkInfo) string {
	mode := info.Mode()
	if mode&os.ModeDir != 0 {
		return f.dirString(info.Name(), info.PathDisplay())
	}
	if mode&os.ModeSymlink != 0 {
		return f.linkString(info, link)
	}
	if mode&os.ModeDevice != 0 {
		name := info.PathDisplay()
		if f.nerdfont {
			return otherIcons["device"] + " " + name
		}
		if f.icons {
			return "💽 " + name
		}
		return name
	}
	if mode&os.ModeNamedPipe != 0 {
		return info.PathDisplay()
	}
	if mode&os.ModeSocket != 0 {
		return info.PathDisplay()
	}
	return f.fileString(info)
}

func (f *FileNameGetterPlain) linkString(info FileInfo, link *LinkInfo) string {
	name := info.PathDisplay()
	if !link.broken && link.isDir {
		if f.nerdfont {
			var linkIcon string
			if link.broken {
				linkIcon = otherIcons["brokenLink"]
			} else {
				linkIcon = otherIcons["linkDir"]
			}
			return linkIcon + " " + name
		} else if f.icons {
			return "🔗 " + name
		}
		return name
	}
	if f.nerdfont {
		if link.broken {
			return otherIcons["brokenLink"] + " " + name
		}
		return otherIcons["link"] + " " + name
	} else if f.icons {
		return "🔗 " + name + " "
	}
	return name
}

func (f *FileNameGetterPlain) linkTargetString(link *LinkInfo) string {
	if link.broken {
		return link.targetDisplay
	}
	if link.isDir {
		return f.dirString(link.target, link.targetDisplay)
	}
	if link.info == nil {
		return link.targetDisplay
	}
	info := link.info
	mode := info.Mode()
	if mode&os.ModeDir != 0 {
		return f.dirString(info.Name(), info.PathDisplay())
	}
	name := info.PathDisplay()
	if mode&os.ModeSymlink != 0 {
		if f.nerdfont {
			if link.broken {
				return otherIcons["brokenLink"] + " " + name
			}
			return otherIcons["link"] + " " + name
		} else if f.icons {
			return "🔗 " + name
		}
		return name + " "
	}
	if mode&os.ModeDevice != 0 {
		if f.nerdfont {
			return otherIcons["device"] + " " + name
		} else if f.icons {
			return "💽 " + name
		}
		return name
	}
	if mode&os.ModeNamedPipe != 0 {
		return name
	}
	if mode&os.ModeSocket != 0 {
		return name
	}
	return f.fileString(info)
}

func (f *FileNameGetterPlain) fileIcon(info FileInfo) string {
	if f.nerdfont {
		if isExecutableFile(info) {
			return getIconForFile("", "shell") + " "
		}
		return getIconForFile(info.Basename(), info.Ext()) + " "
	}
	if f.icons {
		if isExecutableFile(info) {
			return ">_"
		}
	}
	return ""
}

func (f *FileNameGetterPlain) fileString(info FileInfo) string {
	// in some cases files have icons if front
	// if nerd font enabled, then it'll be a file-specific icon, or if its an executable script, a little shell icon
	// if the regular --icons flag is used instead, then it will show a ">_" only if the file is executable
	return f.fileIcon(info) + info.PathDisplay()
}

func (f *FileNameGetterPlain) dirString(name string, pathDisplay string) string {
	icon := ""
	if f.icons {
		icon = "📂 "
	} else if f.nerdfont {
		icon = getIconForFolder(name) + " "
	}
	return icon + pathDisplay
}
