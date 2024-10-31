package csv

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

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type AppInterface interface {
	WorkDir() string
	Abs(path string) (string, error)
	UserHomeDir() (string, error)
	Dir(path string) string
	Join(elem ...string) string
	SplitAll(path string) []string
}

func New(
	app AppInterface,
	args *lsargs.Arguments,
	colors *lscolors.HtmlColors,
) *HtmlFormatter {
	return &HtmlFormatter{
		app:    app,
		args:   args,
		colors: colors,
	}
}

type HtmlFormatter struct {
	app    AppInterface
	args   *lsargs.Arguments
	colors *lscolors.HtmlColors
}

func (*HtmlFormatter) LongLinkTarget() bool {
	return true
}

func (*HtmlFormatter) LinkTargetSep() string {
	return "► "
}

func (*HtmlFormatter) DefaultTimeStyle() string {
	timeStyle := os.Getenv("LSGO_TIME_STYLE")
	if timeStyle != "" {
		return timeStyle
	}
	return "long-iso"
}

func (f *HtmlFormatter) SizeFormat() SizeFormat {
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

func (f *HtmlFormatter) FileError(w io.Writer, errArg error, path string) {
	absPath, err := f.app.Abs(path)
	check(err)
	msg := "► " + absPath
	if f.colors != nil {
		msg = f.Colorize(msg, f.colors.FolderHeader.Error)
	}
	fmt.Fprintln(w, msg)
	fmt.Fprintln(w, errArg.Error())
}

type getPathIface interface {
	GetPath() string
}

func (f *HtmlFormatter) PrintError(w io.Writer, err error) {
	path := ""
	getPath, ok := err.(getPathIface)
	if ok {
		path = getPath.GetPath()
	}
	f.FileError(w, err, path)
}

func (f *HtmlFormatter) prettifyPath(path string) string {
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

func (f *HtmlFormatter) FolderHeader(w io.Writer, path string, itemCount int) {
	if len(*f.args.Find) > 0 && itemCount == 0 {
		return
	}
	if len(f.args.Paths) == 1 && f.args.Paths[0] == "." && !*f.args.Recursive {
		return
	}
	fhColors := f.colors.FolderHeader
	headerString := f.Colorize("►", fhColors.Arrow) + f.Colorize(" ", fhColors.Main)
	prettyPath := f.prettifyPath(path)

	if f.app.Dir(prettyPath) == prettyPath {
		headerString += prettyPath
	} else {
		folders := f.app.SplitAll(prettyPath)
		coloredFolders := make([]string, 0, len(folders))
		for i, folder := range folders {
			if i == len(folders)-1 { // different color for the last folder in the path
				coloredFolders = append(coloredFolders, f.Colorize(folder, fhColors.LastFolder))
				continue
			}
			coloredFolders = append(coloredFolders, f.Colorize(folder, fhColors.Main))
		}
		coloredFolders = append(coloredFolders, "")
		headerString += f.app.Join(coloredFolders...)
		// headerString += f.app.JoinColor(fhColors.Slash.S(), "", coloredFolders...)
	}

	if f.colors != nil {
		headerString += " "
	}
	fmt.Fprintln(w, headerString)
}

func (*HtmlFormatter) FolderTail(w io.Writer, _ string) {
	fmt.Fprintln(w, "<br/>")
}

func (*HtmlFormatter) TableHeader(w io.Writer, tableObj *table.Table) {
	h := []string{}
	for _, col := range tableObj.Columns {
		h = append(h, col.Title)
	}
	fmt.Fprintln(w, strings.Join(h, ","))
}

func (f *HtmlFormatter) Colorize(str string, style *lscolors.Style) string {
	def := f.colors.Default
	fg := def.Fg
	if style.Fg > 0 {
		fg = style.Fg
	}
	bg := def.Bg
	if style.Bg > 0 {
		bg = style.Bg
	}
	css := []string{
		`color:` + lscolors.TermColorsHex[int(fg)],
		`background:` + lscolors.TermColorsHex[int(bg)],
	}
	if style.Bold {
		css = append(css, `font-weight:bold`)
	}
	if strings.HasPrefix(str, " ") {
		str = "&nbsp;" + strings.Trim(str, " ") + "&nbsp;"
	}
	return fmt.Sprintf(`<span style='%s;'>%s</span>`, strings.Join(css, ";"), str)
}

func (f *HtmlFormatter) FormatValue(_ string, value any) (string, error) {
	// _: colName
	switch valueTyped := value.(type) {
	case string:
		return valueTyped, nil
	case uint64:
		return strconv.FormatUint(valueTyped, 10), nil
	}
	return fmt.Sprintf("%v", value), nil
}

func (f *HtmlFormatter) FormatItem(tableObj *table.Table, item any) ([]string, error) {
	return tableObj.FormatItem(item)
}

func (f *HtmlFormatter) PrintItems(w io.Writer, _ *table.Table, items iface.FormattedItemList) error {
	bgColor := lscolors.TermColorsHex[int(f.colors.Default.Bg)]
	fmt.Fprintf(w, `<table style="background-color:%s;font-family: monospace;">\n`, bgColor)
	for i := 0; i < items.Len(); i++ {
		row := items.Get(i)
		tdList := make([]string, len(row))
		for i, cell := range row {
			tdList[i] = "<td>" + cell + "</td>"
		}
		fmt.Fprintln(
			w,
			"<tr>"+strings.Join(tdList, "")+"</tr>",
		)
	}
	fmt.Fprintln(w, "</table>")
	return nil
}
