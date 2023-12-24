package application

import (
	"log"
	"os"

	"github.com/ilius/go-table"
	c "github.com/ilius/ls-go/common"
	"github.com/ilius/ls-go/filesystem"
	lscsv "github.com/ilius/ls-go/format/csv"
	lshtml "github.com/ilius/ls-go/format/html"
	lsjson "github.com/ilius/ls-go/format/json"
	jsonarray "github.com/ilius/ls-go/format/jsonarray"
	"github.com/ilius/ls-go/format/tabular"
	"github.com/ilius/ls-go/iface"
	"github.com/ilius/ls-go/lsargs"
	"github.com/ilius/ls-go/lscolors"
	"github.com/ilius/ls-go/lsplatform"
	"github.com/ilius/ls-go/lstime"
	"github.com/ilius/ls-go/terminal"
)

const VERSION = "1.2.0"

var getOwnerAndGroup func(lsplatform.FileInfo) (*lsplatform.OwnerGroup, error)

var (
	platform     = lsplatform.New()
	rootUserName = platform.RootUserName()
)

type Application struct {
	iface.FileSystem
	iface.Terminal
	iface.Formatter

	Platform *lsplatform.LocalPlatform
	workDir  string

	PrimaryTimeColName string

	QuotingStyle string
	EnsureASCII  bool

	exitStatus int

	errors []error

	QuestionMark string
}

func NewApplication() *Application {
	fs := filesystem.NewLocalFileSystem()
	return &Application{
		FileSystem: fs,
		Platform:   platform,
		Terminal:   terminal.NewLocalTerminal(),
		workDir:    fs.WorkDir(),
	}
}

func (app *Application) AddError(err error) {
	app.errors = append(app.errors, err)
}

func (app *Application) PrintErrors() {
	f := app.Formatter
	for _, err := range app.errors {
		f.PrintError(stderr, err)
	}
}

func (app *Application) Exit() {
	if app.exitStatus != 0 {
		os.Exit(app.exitStatus)
	}
}

func (app *Application) onFileError(err error, path string) {
	if os.IsNotExist(err) || os.IsPermission(err) {
		app.exitStatus = 2
		app.AddError(&c.FileError{
			Path: path,
			Msg:  err.Error(),
		})
		return
	}
	panic(err)
}

func (app *Application) WorkDir() string {
	return app.workDir
}

func (app *Application) makeFormatter(colorsEnabled bool) iface.Formatter {
	if *args.Json {
		return lsjson.New(args)
	}
	if *args.JsonArray {
		return jsonarray.New(args)
	}
	if *args.Csv {
		return lscsv.New(args)
	}
	if *args.Html {
		return lshtml.New(app, args, colors.Html)
	}
	if colorsEnabled {
		return tabular.New(app, args, colors.Tabular)
	}
	return tabular.New(app, args, nil)
}

func (app *Application) longSet(cols map[string]bool, nameParams *FileNameParams) {
	cols[c.C_Mode] = true
	cols[c.C_HardLinks] = true
	cols[c.C_Owner] = true
	cols[c.C_Group] = true
	cols[c.C_Size] = true
	cols[app.PrimaryTimeColName] = true
	if app.Formatter.LongLinkTarget() {
		cols[c.C_LinkTarget] = true
		nameParams.showLinks = true
	}
}

func (app *Application) extraLongSet(cols map[string]bool, nameParams *FileNameParams) {
	app.longSet(cols, nameParams)

	cols[c.C_Inode] = true
	cols[c.C_ModeOct] = true
	cols[c.C_Blocks] = true
	cols[c.C_MTime] = true
	cols[c.C_CTime] = true
	cols[c.C_ATime] = true
}

func (app *Application) PostParse(args *lsargs.Arguments) *table.TableSpec {
	colors, err := app.Terminal.ColorsEnabled(*args.Color)
	check(err)

	formatter := app.makeFormatter(colors)
	app.Formatter = formatter

	app.QuestionMark = formatter.Colorize("?", lscolors.Fg(1))

	cols := map[string]bool{}

	if *args.FullTime {
		*args.Long = true
		*args.TimeStyle = "full-iso"
	}
	if *args.NumericUidGid {
		*args.Long = true
		getOwnerAndGroup = app.Platform.OwnerAndGroupIDs
	} else {
		getOwnerAndGroup = app.Platform.OwnerAndGroupNames
	}
	if *args.DirsOnly && *args.FilesOnly {
		log.Fatal("--dirs-only and --files cannot both be set")
	}
	if *args.Nerdfont && *args.Icons {
		log.Fatal("--nerd-font and --icons cannot both be set")
	}

	if *args.Shortcut_t {
		*args.Sort = c.S_TIME
	}
	if *args.Shortcut_U {
		*args.Sort = c.S_NONE
	}
	if *args.Shortcut_S {
		*args.Sort = c.S_SIZE
	}
	if *args.Shortcut_X {
		*args.Sort = c.S_EXTENSION
	}

	{
		timeCol := *args.Time
		if *args.Shortcut_c {
			timeCol = "ctime"
		} else if *args.Shortcut_u {
			timeCol = "use"
		}
		app.PrimaryTimeColName = timeColumnFromInput(timeCol)
	}

	quotingStyle := os.Getenv("QUOTING_STYLE")
	if *args.QuotingStyle != "" {
		quotingStyle = *args.QuotingStyle
	}
	if *args.Shortcut_literal {
		quotingStyle = c.E_literal
	}
	if *args.Shortcut_escape {
		quotingStyle = c.E_escape
	}
	switch quotingStyle {
	case "":
		quotingStyle = c.E_shell_escape
	case c.E_locale:
		log.Fatalf("unsupported --quoting-style=locale")
	case c.E_none, c.E_literal, c.E_shell, c.E_shell_always, c.E_shell_escape, c.E_shell_escape_always, c.E_c, c.E_escape:
		break
	default:
		log.Fatalf("invalid --quoting-style=%v", quotingStyle)
	}
	app.QuotingStyle = quotingStyle

	app.EnsureASCII = *args.ASCII

	nameParams := &FileNameParams{
		showLinksSep: formatter.LinkTargetSep(),
		showLinks:    *args.Links,
		linkRel:      *args.LinkRel,
		icons:        *args.Icons,
		nerdfont:     *args.Nerdfont,
		fullPath:     len(args.Paths) > 1 || *args.Recursive,
	}

	if *args.Long {
		app.longSet(cols, nameParams)
	}
	if *args.ExtraLong {
		app.extraLongSet(cols, nameParams)
	}
	if *args.Inode {
		cols[c.C_Inode] = true
	}
	if *args.Blocks {
		cols[c.C_Blocks] = true
	}
	if *args.ModeOct {
		cols[c.C_ModeOct] = true
	}
	if *args.Mode {
		cols[c.C_Mode] = true
	}
	if *args.Owner {
		cols[c.C_Owner] = true
	}
	if *args.Group {
		cols[c.C_Group] = true
	}
	if *args.NoGroup {
		cols[c.C_Group] = false
	}
	if *args.Size {
		cols[c.C_Size] = true
	}
	if *args.Mtime {
		cols[c.C_MTime] = true
	}
	if *args.Ctime {
		cols[c.C_CTime] = true
	}
	if *args.Atime {
		cols[c.C_ATime] = true
	}
	cols[c.C_Name] = true

	timeParams := &lstime.TimeParams{}
	timeStyle := formatter.DefaultTimeStyle()
	if *args.TimeStyle != "" {
		timeStyle = *args.TimeStyle
	}
	check(timeParams.SetTimeStyle(timeStyle))

	exprList := []string{}
	if *args.Expr != "" {
		exprList = append(exprList, *args.Expr)
	}

	tableSpec := makeTableSpec(
		cols,
		formatter,
		colors,
		nameParams,
		timeParams,
		exprList,
	)
	return tableSpec
}
