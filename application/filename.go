package application

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/ilius/go-table"
	. "github.com/ilius/ls-go/common"
	"github.com/ilius/ls-go/lscolors"
)

func NewFileNameGetter(colors bool, params *FileNameParams) table.Getter {
	if colors {
		return &FileNameGetter{params}
	}
	return &FileNameGetterPlain{params}
}

type FileNameParams struct {
	showLinksSep string
	fullPath     bool
	showLinks    bool
	linkRel      bool
	icons        bool
	nerdfont     bool
}

// check for executable permissions
func isExecutableFile(info FileInfo) bool {
	m := info.Mode()
	return m&0o111 != 0 && m&fs.ModeType == 0
}

type FileNameGetter struct {
	*FileNameParams
}

func (f *FileNameGetter) Value(item any) (any, error) {
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
		displayName += app.Colorize("► ", colors.Link.Arrow) + f.linkTargetString(link)
	}

	return displayName, nil
}

func (f *FileNameGetter) ValueString(colName string, item any) (string, error) {
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

func (f *FileNameGetter) Format(item any, value any) (string, error) {
	return value.(string), nil
}

func (f *FileNameGetter) nameString(info FileInfo, link *LinkInfo) string {
	mode := info.Mode()
	if mode&os.ModeDir != 0 {
		return f.dirString(info.Name(), info.PathDisplay())
	}
	if mode&os.ModeSymlink != 0 {
		return f.linkString(info, link)
	}
	if mode&os.ModeDevice != 0 {
		name := info.PathDisplay()
		color := colors.Device
		if f.nerdfont {
			return app.Colorize(otherIcons["device"]+" "+name+" ", color)
		}
		if f.icons {
			return app.Colorize("💽 "+name+" ", color)
		}
		return app.Colorize(" "+name+" ", color)
	}
	if mode&os.ModeNamedPipe != 0 {
		return app.Colorize(" "+info.PathDisplay()+" ", colors.Pipe)
	}
	if mode&os.ModeSocket != 0 {
		return app.Colorize(" "+info.PathDisplay()+" ", colors.Socket)
	}
	return f.fileString(info)
}

func (f *FileNameGetter) linkString(info FileInfo, link *LinkInfo) string {
	name := info.PathDisplay()
	if !link.broken && link.isDir {
		color := colors.Link.NameDir
		if f.nerdfont {
			var linkIcon string
			if link.broken {
				linkIcon = otherIcons["brokenLink"]
			} else {
				linkIcon = otherIcons["linkDir"]
			}
			return app.Colorize(linkIcon+" "+name+" ", color) + " "
		} else if f.icons {
			return app.Colorize("🔗 "+name+" ", color) + " "
		}
		return app.Colorize(" "+name+" ", color) + " "
	}
	color := colors.Link.Name
	if f.nerdfont {
		if link.broken {
			return app.Colorize(otherIcons["brokenLink"]+" "+name+" ", color)
		}
		return app.Colorize(otherIcons["link"]+" "+name+" ", color)
	}
	if f.icons {
		return app.Colorize("🔗 "+name+" ", color)
	}
	return app.Colorize(name+" ", color)
}

func (f *FileNameGetter) linkTargetString(link *LinkInfo) string {
	if link.broken {
		return app.Colorize(link.targetDisplay, colors.Link.Broken)
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
		color := colors.Link.Name
		if f.nerdfont {
			if link.broken {
				return app.Colorize(otherIcons["brokenLink"]+" "+name+" ", color)
			}
			return app.Colorize(otherIcons["link"]+" "+name+" ", color)
		} else if f.icons {
			return app.Colorize("🔗 "+name+" ", color)
		}
		return app.Colorize(name+" ", color)
	}
	if mode&os.ModeDevice != 0 {
		color := colors.Device
		if f.nerdfont {
			return app.Colorize(otherIcons["device"]+" "+name+" ", color)
		}
		if f.icons {
			return app.Colorize("💽 "+name+" ", color)
		}
		return app.Colorize(" "+name+" ", color)
	}
	if mode&os.ModeNamedPipe != 0 {
		return app.Colorize(" "+name+" ", colors.Pipe)
	}
	if mode&os.ModeSocket != 0 {
		return app.Colorize(" "+name+" ", colors.Socket)
	}
	return f.fileString(info)
}

func (f *FileNameGetter) fileIcon(info FileInfo, mainColor *lscolors.Style) string {
	if f.nerdfont {
		if isExecutableFile(info) {
			return app.Colorize(getIconForFile("", "shell")+" ", mainColor)
		}
		return app.Colorize(getIconForFile(info.Basename(), info.Ext())+" ", mainColor)
	} else if f.icons {
		if isExecutableFile(info) {
			return app.Colorize(">_", lscolors.BgGray(1).SetFg(46)) + " "
		}
	}
	return ""
}

func (f *FileNameGetter) fileString(info FileInfo) string {
	basename := info.Basename()
	ext := info.Ext()
	suffix := info.Suffix()
	key := ""
	if ext != "" {
		key = strings.ToLower(ext)[1:]
	}
	// figure out which color to choose
	color := colors.File[lscolors.DEFAULT]
	alias, hasAlias := FileAliases[key]
	if hasAlias {
		key = alias
	}
	betterColor, hasBetterColor := colors.File[key]
	if hasBetterColor {
		color = betterColor
	}
	mainColor, accentColor := color.Get()

	// in some cases files have icons if front
	// if nerd font enabled, then it'll be a file-specific icon, or if its an executable script, a little shell icon
	// if the regular --icons flag is used instead, then it will show a ">_" only if the file is executable
	icon := f.fileIcon(info, mainColor)
	quotedBasename := quoteFileName(basename)
	colorize := app.Colorize
	if quotedBasename != basename {
		return icon + colorize(info.PathDisplay(), mainColor)
	}
	return strings.Join([]string{
		icon,
		colorize(basename, mainColor),
		colorize(ext, accentColor),
		suffix,
	}, "")
}

func (f *FileNameGetter) dirString(name string, pathDisplay string) string {
	color := colors.Dir.Name
	if strings.HasPrefix(name, ".") {
		color = colors.Dir.HiddenName
	}
	icon := " "
	if f.icons {
		icon = "📂 "
	} else if f.nerdfont {
		icon = getIconForFolder(name) + " "
	}
	return app.Colorize(
		icon+pathDisplay+" ",
		color,
	)
}
