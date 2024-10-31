package json

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/ilius/go-table"
	. "github.com/ilius/ls-go/common"
	"github.com/ilius/ls-go/iface"
	"github.com/ilius/ls-go/lsargs"
	"github.com/ilius/ls-go/lscolors"
)

func New(args *lsargs.Arguments) *JsonFormatter {
	return &JsonFormatter{
		args:         args,
		ensure_ascii: *args.ASCII,
	}
}

type JsonFormatter struct {
	args *lsargs.Arguments

	ensure_ascii bool
}

func (*JsonFormatter) LongLinkTarget() bool {
	return true
}

func (*JsonFormatter) LinkTargetSep() string {
	return ","
}

func (*JsonFormatter) DefaultTimeStyle() string {
	return "full-iso"
}

func (*JsonFormatter) SizeFormat() SizeFormat {
	return SizeFormatInteger
}

func (*JsonFormatter) FileError(w io.Writer, err error, path string) {
	jsonB, encodeErr := json.Marshal(ItemErrorJSON{
		Name:  path,
		Error: err.Error(),
	})
	fmt.Fprintln(w, string(jsonB))
	if encodeErr != nil {
		panic(encodeErr)
	}
}

func (f *JsonFormatter) PrintError(w io.Writer, err error) {
	// TODO: check if its file error, then call FileError
	b, jerr := json.Marshal(err)
	if jerr == nil {
		_, _ = w.Write(b)
		_, _ = w.Write([]byte{'\n'})
		return
	}
	b, _ = json.Marshal(map[string]string{
		"error": err.Error(),
	})
	_, _ = w.Write(b)
	_, _ = w.Write([]byte{'\n'})
}

func (*JsonFormatter) FolderHeader(_ io.Writer, _ string, _ int) {}

func (*JsonFormatter) FolderTail(_ io.Writer, _ string) {}

func (f *JsonFormatter) TableHeader(w io.Writer, tableObj *table.Table) {
	if !*f.args.Header {
		return
	}
	parts := []string{}
	for _, col := range tableObj.Columns {
		part, err := jsonKeyValue(col.Name, col.Title, false)
		if err != nil {
			panic(err)
		}
		parts = append(parts, part)
	}
	fmt.Fprintln(w, "{"+strings.Join(parts, ",")+"}")
}

func (f *JsonFormatter) Colorize(str string, _ *lscolors.Style) string {
	return str
}

func (f *JsonFormatter) FormatValue(colName string, value any) (string, error) {
	return jsonKeyValue(colName, value, f.ensure_ascii)
}

func (*JsonFormatter) FormatItem(tableObj *table.Table, item any) ([]string, error) {
	str, err := tableObj.FormatItemBasic(item, ",")
	if err != nil {
		return nil, err
	}
	return []string{
		"{" + str + "}",
	}, nil
}

func (*JsonFormatter) PrintItems(w io.Writer, _ *table.Table, items iface.FormattedItemList) error {
	for i := 0; i < items.Len(); i++ {
		fmt.Fprintln(w, items.Get(i)[0])
	}
	return nil
}
