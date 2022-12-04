package application

import "os"

// LinkInfo wraps link stat info and whether the link points to valid file
type LinkInfo struct {
	info          FileInfo
	target        string
	targetDisplay string
	isDir         bool
	broken        bool
	relative      bool
}

func (link *LinkInfo) Dereference(source FileInfo) FileInfo {
	infoOrig := link.info
	if infoOrig == nil {
		return nil
	}
	info := NewFileInfoLowFrom(infoOrig)
	name := source.Name()
	info.name = name
	pname := app.FileSystem.SplitExt(name)
	return &FileInfoImp{
		FileInfo: info,
		dir:      source.Dir(),
		curDir:   source.CurDir(),
		isAbs:    source.IsAbs(),
		basename: pname.Base,
		ext:      pname.Ext,
		suffix:   pname.Suffix,
	}
}

func getLinkInfo(info FileInfo, parentDirAbs string, rel bool) *LinkInfo {
	absPath := app.FileSystem.Join(parentDirAbs, info.Name())
	target, err1 := app.FileSystem.ReadLink(absPath)
	check(err1)

	targetAbs := target
	targetDisplay := target
	relative := false
	if !app.FileSystem.IsAbs(target) {
		targetAbs = app.FileSystem.Join(parentDirAbs, target)
		relative = true
	}
	linkInfo, err2 := app.FileSystem.Stat(targetAbs)
	if rel {
		linkRel, _ := app.FileSystem.Rel(parentDirAbs, target)
		if linkRel != "" && len(linkRel) <= len(target) {
			// i prefer the look of these relative paths prepended with ./
			if linkRel[0] != '.' {
				targetDisplay = "./" + linkRel
			} else {
				targetDisplay = linkRel
			}
		}
	}

	targetDisplay = quoteFileName(targetDisplay)

	link := &LinkInfo{
		target:        target,
		targetDisplay: targetDisplay,
		relative:      relative,
	}
	if linkInfo != nil {
		link.isDir = linkInfo.IsDir()
		targetParsed := app.FileSystem.SplitExt(target)
		link.info = &FileInfoImp{
			FileInfo: linkInfo,
			basename: targetParsed.Base,
			ext:      targetParsed.Ext,
			suffix:   targetParsed.Suffix,
			dir:      app.FileSystem.Dir(target),
			curDir:   parentDirAbs,
			isAbs:    !relative,
		}
	} else if os.IsNotExist(err2) {
		link.broken = true
	} else if !os.IsPermission(err2) {
		check(err2)
	}
	return link
}
