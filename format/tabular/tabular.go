package tabular

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/ilius/go-table"
	. "github.com/ilius/ls-go/common"
	"github.com/ilius/ls-go/iface"
	"github.com/ilius/ls-go/lsargs"
	"github.com/ilius/ls-go/lscolors"
)

const Reset = "\x1b[0m"

type AppInterface interface {
	WorkDir() string
	TermWidth() (int, error)
	Abs(path string) (string, error)
	UserHomeDir() (string, error)
	Dir(path string) string
	Join(elem ...string) string
	SplitAll(path string) []string
	JoinColor(color string, reset string, elem ...string) string
}

func New(
	app AppInterface,
	args *lsargs.Arguments,
	colors *lscolors.TabularColors,
) *TabularFormatter {
	return &TabularFormatter{
		app:    app,
		args:   args,
		colors: colors,
	}
}

type TabularFormatter struct {
	app    AppInterface
	args   *lsargs.Arguments
	colors *lscolors.TabularColors
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func (*TabularFormatter) LongLinkTarget() bool {
	return true
}

func (*TabularFormatter) LinkTargetSep() string {
	return "► "
}

func (*TabularFormatter) DefaultTimeStyle() string {
	timeStyle := os.Getenv("LSGO_TIME_STYLE")
	if timeStyle != "" {
		return timeStyle
	}
	return "long-iso"
}

func (f *TabularFormatter) SizeFormat() SizeFormat {
	// --si overrides --human-readable
	if *f.args.SI {
		// display human-readable format for size, but with powers of 1000
		return SizeFormatMetric
	}
	if *f.args.Human {
		return SizeFormatLegacy // human-readable
	}
	if *f.args.Bytes {
		return SizeFormatInteger
	}
	return SizeFormatLegacy // human-readable
}

func (f *TabularFormatter) prettifyPath(path string) string {
	prettyPath, err := f.app.Abs(path)
	check(err)
	home, err := f.app.UserHomeDir()
	check(err)

	// FIXME: this is buggy and it sucks
	workDir := f.app.WorkDir()
	if strings.HasPrefix(prettyPath, workDir) {
		prettyPath = "." + prettyPath[len(workDir):]
	} else if strings.HasPrefix(prettyPath, home) {
		prettyPath = "~" + prettyPath[len(home):]
	}
	return prettyPath
}

func (f *TabularFormatter) FileError(w io.Writer, err error, path string) {
	msg := "► " + f.prettifyPath(path)
	if f.colors != nil {
		msg = f.Colorize(msg, f.colors.FolderHeader.Error)
	}
	fmt.Fprintln(w, msg)
	fmt.Fprintln(w, err.Error())
}

func (f *TabularFormatter) PrintError(w io.Writer, err error) {
	fileErr, ok := err.(*FileError)
	if ok {
		f.FileError(w, err, fileErr.Path)
		return
	}
	fmt.Fprintln(w, err.Error())
}

func (f *TabularFormatter) FolderHeader(w io.Writer, path string, itemCount int) {
	if len(*f.args.Find) > 0 && itemCount == 0 {
		return
	}
	if len(f.args.Paths) == 1 && f.args.Paths[0] == "." && !*f.args.Recursive {
		return
	}
	if f.colors != nil {
		f.folderHeader(w, path, itemCount)
		return
	}
	f.folderHeaderNoColor(w, path, itemCount)
}

func (f *TabularFormatter) folderHeader(w io.Writer, path string, _ int) {
	// _: itemCount
	fhColors := f.colors.FolderHeader
	headerString := fhColors.Arrow.S() + "►" + fhColors.Main.S() + " "
	prettyPath := f.prettifyPath(path)

	if f.app.Dir(prettyPath) == prettyPath {
		headerString += prettyPath
	} else {
		folders := f.app.SplitAll(prettyPath)
		coloredFolders := make([]string, 0, len(folders))
		for i, folder := range folders {
			if i == len(folders)-1 { // different color for the last folder in the path
				coloredFolders = append(coloredFolders, fhColors.LastFolder.S()+folder)
			} else {
				coloredFolders = append(coloredFolders, fhColors.Main.S()+folder)
			}
		}
		coloredFolders = append(coloredFolders, "")
		headerString += f.app.JoinColor(fhColors.Slash.S(), "", coloredFolders...)
	}

	if f.colors != nil {
		headerString += " " + Reset
	}
	fmt.Fprintln(w, headerString)
}

func (f *TabularFormatter) folderHeaderNoColor(w io.Writer, path string, _ int) {
	// _: itemCount
	headerString := "► "
	prettyPath := f.prettifyPath(path)

	if f.app.Dir(prettyPath) == prettyPath {
		headerString += prettyPath
	} else {
		folders := f.app.SplitAll(prettyPath)
		folders = append(folders, "")
		headerString += f.app.Join(folders...)
	}

	if f.colors != nil {
		headerString += " " + Reset
	}
	fmt.Fprintln(w, headerString)
}

func (*TabularFormatter) FolderTail(w io.Writer, _ string) {
	fmt.Fprintln(w, "")
}

func (f *TabularFormatter) getSep() string {
	if *f.args.Vbar {
		return " | "
	}
	return " "
}

func (f *TabularFormatter) TableHeader(w io.Writer, tableObj *table.Table) {
	if !*f.args.Header {
		return
	}
	if !f.oneFilePerLine(tableObj) {
		return
	}
	sep := f.getSep()
	if f.colors != nil {
		fmt.Fprint(w, f.Colorize(tableObj.FormatHeader(sep), f.colors.TableHeader))
		return
	}
	fmt.Fprint(w, tableObj.FormatHeader(sep))
}

func (f *TabularFormatter) Colorize(str string, style *lscolors.Style) string {
	if style == nil {
		return str
	}
	return style.S() + str + Reset
}

func (*TabularFormatter) FormatValue(_ string, value any) (string, error) {
	// _: colName
	switch valueTyped := value.(type) {
	case string:
		return valueTyped, nil
	case uint64:
		return strconv.FormatUint(valueTyped, 10), nil
	}
	return fmt.Sprintf("%v", value), nil
}

func (*TabularFormatter) FormatItem(tableObj *table.Table, item any) ([]string, error) {
	return tableObj.FormatItem(item)
}

// return true to print one file per line
// return false to print multiple files in a line
func (f *TabularFormatter) oneFilePerLine(tableObj *table.Table) bool {
	if *f.args.Long {
		return true
	}
	if *f.args.ExtraLong {
		return true
	}
	if *f.args.SingleCol {
		return true
	}
	if *f.args.Horizontal {
		return false
	}
	if *f.args.Vertical {
		return false
	}
	if len(tableObj.Columns) > 3 {
		// too many columns, many files per line makes it confusing
		return true
	}
	return false
}

func (f *TabularFormatter) PrintItems(w io.Writer, tableObj *table.Table, items iface.FormattedItemList) error {
	count := items.Len()
	sep := f.getSep()
	if f.oneFilePerLine(tableObj) {
		for index := 0; index < count; index++ {
			aligned, err := tableObj.AlignFormattedItem(items.Get(index))
			if err != nil {
				return err
			}
			fmt.Fprintln(w, strings.Join(aligned, sep))
		}
		return nil
	}

	// format in columns, like `ls` or `ls -x`
	maxWidth, err := f.app.TermWidth()
	if err != nil {
		return err
	}
	if *f.args.Horizontal {
		tableObj.MergeRowsHorizontal(w, items, maxWidth, sep, *f.args.Compact)
	} else {
		tableObj.MergeRowsVertical(w, items, maxWidth, sep, *f.args.Compact)
	}
	return nil
}
