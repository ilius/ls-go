package application

import (
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"strconv"

	"github.com/ilius/go-table"
	jsonparse "github.com/ilius/ls-go/parse/json"
)

func (app *Application) ListMain(tableSpec *table.TableSpec) {
	if *args.ReadJson {
		res, err := jsonparse.Parse(os.Stdin)
		check(err)
		app.ListFiles(
			table.NewTable(tableSpec),
			app.WorkDir(),
			res.Files,
			true, // forceDotfiles
		)
		return
	}

	// separate the directories to be listed from other files
	dirs := []string{}
	files := []FileInfo{}
	dirListEnable := !*args.Directory
	for _, path := range args.Paths {
		fileStat, err := app.FileSystem.Stat(path)
		if err != nil {
			app.onFileError(err, path)
			continue
		}
		if fileStat.IsDir() && dirListEnable {
			dirs = append(dirs, path)
			continue
		}
		pname := app.FileSystem.SplitExt(path)
		files = append(files, &FileInfoImp{
			FileInfo: fileStat,
			basename: pname.Base,
			ext:      pname.Ext,
			suffix:   pname.Suffix,
			dir:      app.FileSystem.Dir(path),
			isAbs:    app.FileSystem.IsAbs(path),
		})
	}

	// list files first
	if len(files) > 0 {
		app.ListFiles(
			table.NewTable(tableSpec),
			app.WorkDir(),
			files,
			true, // forceDotfiles
		)
		if len(dirs) > 0 {
			app.FolderTail(stdout, "")
		}
	}

	// then list the contents of each directory
	app.ListDirList(tableSpec, dirs)
}

func (app *Application) ListDirList(tableSpec *table.TableSpec, pathList []string) {
	if len(pathList) == 0 {
		return
	}
	count := len(pathList)
	for _, path := range pathList[:count-1] {
		count := app.ListDir(table.NewTable(tableSpec), path)
		if count > 0 {
			app.FolderTail(stdout, path)
		}
	}
	app.ListDir(table.NewTable(tableSpec), pathList[count-1])
}

func (app *Application) ListDir(tableObj *table.Table, path string) int {
	items := []FileInfo{}

	pathAbs, err := app.FileSystem.Abs(path)
	check(err)

	addItem := func(info *fs.FileInfo) {
		pname := app.FileSystem.SplitExt((*info).Name())
		items = append(items, &FileInfoImp{
			FileInfo: *info,
			basename: pname.Base,
			ext:      pname.Ext,
			suffix:   pname.Suffix,
			dir:      pathAbs,
			curDir:   pathAbs,
			isAbs:    false,
		})
	}

	entries, err := app.FileSystem.ReadDir(path)
	// if we couldn't read the folder, print a "header" with error message and use error-looking colors
	if err != nil {
		app.onFileError(err, path)
		return 0
	}

	for _, entry := range entries {
		info, err := entry.Info()
		check(err)
		addItem(&info)
	}
	if *args.All {
		for _, name := range []string{".", ".."} {
			info, err := app.FileSystem.Stat(name)
			if err != nil {
				app.onFileError(err, name)
				return 0
			}
			addItem(&info)
		}
	}

	// filter by the regexp if one was passed
	if len(*args.Find) > 0 {
		re, err := regexp.Compile(*args.Find)
		check(err)
		filteredItems := []FileInfo{}
		for _, fileInfo := range items {
			if re.MatchString(fileInfo.Name()) {
				filteredItems = append(filteredItems, fileInfo)
			}
		}
		items = filteredItems
	}

	app.FolderHeader(stdout, path, len(items))

	if len(items) > 0 {
		app.ListFiles(tableObj, path, items, false)
	}

	count := len(items)

	if *args.Recursive {
		for _, item := range items {
			if item.IsDir() && (item.Name()[0] != '.' || *args.All || *args.AlmostAll) {
				if count > 0 {
					app.FolderTail(stdout, item.Name())
				}
				count = app.ListDir(tableObj, app.FileSystem.Join(path, item.Name()))
			}
		}
	}

	return count
}

func (app *Application) ListFiles(tableObj *table.Table, parentDir string, infoList []FileInfo, forceDotfiles bool) {
	if *args.All || *args.AlmostAll {
		forceDotfiles = true
	}

	filesOnly := *args.FilesOnly
	dirsOnly := *args.DirsOnly
	dirsFirst := *args.DirsFirst
	linkRel := *args.LinkRel
	dereference := *args.Dereference

	// collect all the contents here
	files := []*DisplayItem{}
	pinDirs := []*DisplayItem{}

	renderItem := func(info FileInfo) *DisplayItem {
		display, err := app.FormatItem(tableObj, info)
		check(err)
		return &DisplayItem{
			FileInfo: info,
			Time:     info.Time(app.PrimaryTimeColName),
			Display:  display,
		}
	}
	addCondCount := 0
	add := func(info FileInfo) {
		files = append(files, renderItem(info))
	}
	addFile := func(info FileInfo) {
		add(info)
	}
	addDir := func(info FileInfo) {
		add(info)
	}
	if *args.Minsize > 0 {
		minsize := int64(*args.Minsize)
		if *args.Maxsize > 0 {
			maxsize := int64(*args.Maxsize)
			add = func(info FileInfo) {
				if info.Size() >= minsize && info.Size() <= maxsize {
					files = append(files, renderItem(info))
				}
			}
		} else {
			add = func(info FileInfo) {
				if info.Size() >= minsize {
					files = append(files, renderItem(info))
				}
			}
		}
		addCondCount++
	} else if *args.Maxsize > 0 {
		maxsize := int64(*args.Maxsize)
		add = func(info FileInfo) {
			if info.Size() <= maxsize {
				files = append(files, renderItem(info))
			}
		}
		addCondCount++
	}
	if *args.Where != "" {
		getter := NewExprGetter(false, *args.Where)
		if addCondCount == 0 {
			add = func(info FileInfo) {
				if getter.MustValueBool(info) {
					files = append(files, renderItem(info))
				}
			}
		} else {
			lastAdd := add
			add = func(info FileInfo) {
				if getter.MustValueBool(info) {
					lastAdd(info)
				}
			}
		}
		// addCondCount++
	}
	if *args.HasMode != "" {
		mode, err := strconv.ParseUint(*args.HasMode, 8, 16)
		if err != nil {
			panic(fmt.Errorf("--has-mode: bad octal mode %#v", *args.HasMode))
		}
		if addCondCount == 0 {
			add = func(info FileInfo) {
				if uint64(info.Mode())&mode == mode {
					files = append(files, renderItem(info))
				}
			}
		} else {
			lastAdd := add
			add = func(info FileInfo) {
				if uint64(info.Mode())&mode == mode {
					lastAdd(info)
				}
			}
		}
		// addCondCount++
	}
	if filesOnly {
		addDir = func(FileInfo) {}
	} else if dirsFirst {
		addDir = func(info FileInfo) {
			pinDirs = append(pinDirs, renderItem(info))
		}
	}
	if dirsOnly {
		addFile = func(FileInfo) {}
	}
	addSymLink := func(info FileInfo) {
		// read some info about linked file if this info is a symlink
		link := getLinkInfo(info, info.DirAbs(), linkRel)
		if link.info == nil {
			add(info)
			return
		}
		isDir := link.info.IsDir()
		// still not sure if filesOnly and dirsOnly should apply to link targets
		if isDir {
			if filesOnly {
				return
			}
		} else if dirsOnly {
			return
		}
		if dereference {
			info = link.Dereference(info)
		}
		if isDir {
			addDir(info)
			return
		}
		addFile(info)
	}

	for _, info := range infoList {
		// if this is a dotfile (hidden file), we can skip everything with this
		// file if we aren't using the `all` option
		if info.Name()[0] == '.' && !forceDotfiles {
			continue
		}
		if info.Mode()&os.ModeSymlink != 0 {
			addSymLink(info)
			continue
		}
		if info.IsDir() {
			addDir(info)
			continue
		}
		addFile(info)
	}

	// print header after formatting/rendering all items (tableObj.FormatItem)
	// so that we know the width of every column

	app.TableHeader(stdout, tableObj)

	sortFiles(files, *args.Sort, *args.Reverse)
	sortDirs(pinDirs, *args.Sort, *args.Reverse)

	// combine the items together again after sorting, then format and print
	check(app.PrintItems(
		stdout,
		tableObj,
		DisplayItemList(append(pinDirs, files...)),
	))

	if *args.Stats {
		colorsEnable, err := app.Terminal.ColorsEnabled(*args.Color)
		check(err)
		printStats(colorsEnable, len(files), len(pinDirs))
	}
}
