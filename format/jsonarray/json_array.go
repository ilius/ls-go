package jsonarray

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/ilius/go-table"
	. "github.com/ilius/ls-go/common"
	"github.com/ilius/ls-go/escape"
	"github.com/ilius/ls-go/iface"
	"github.com/ilius/ls-go/lsargs"
	"github.com/ilius/ls-go/lscolors"
)

func New(args *lsargs.Arguments) *JsonArrayFormatter {
	return &JsonArrayFormatter{
		args:         args,
		ensure_ascii: *args.ASCII,
	}
}

type JsonArrayFormatter struct {
	args *lsargs.Arguments

	ensure_ascii bool
}

func (*JsonArrayFormatter) LongLinkTarget() bool {
	return true
}

func (*JsonArrayFormatter) LinkTargetSep() string {
	return ","
}

func (*JsonArrayFormatter) DefaultTimeStyle() string {
	return "full-iso"
}

func (*JsonArrayFormatter) SizeFormat() SizeFormat {
	return SizeFormatInteger
}

func (*JsonArrayFormatter) FileError(w io.Writer, err error, path string) {
	jsonB, encodeErr := json.Marshal([]string{
		"error",
		err.Error(),
	})
	fmt.Fprintln(w, string(jsonB))
	if encodeErr != nil {
		panic(encodeErr)
	}
}

func (f *JsonArrayFormatter) PrintError(w io.Writer, err error) {
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

func (*JsonArrayFormatter) FolderHeader(w io.Writer, path string, itemCount int) {}

func (*JsonArrayFormatter) FolderTail(w io.Writer, path string) {}

func (f *JsonArrayFormatter) TableHeader(w io.Writer, tableObj *table.Table) {
	if *f.args.NoHeader {
		return
	}
	row := []string{}
	for _, col := range tableObj.Columns {
		row = append(row, col.Title)
	}
	jsonB, err := json.Marshal(row)
	if err != nil {
		panic(err)
	}
	_, _ = w.Write(jsonB)
	_, _ = w.Write([]byte{'\n'})
}

func (f *JsonArrayFormatter) Colorize(str string, style *lscolors.Style) string {
	return str
}

func (f *JsonArrayFormatter) FormatValue(colName string, value any) (string, error) {
	j_value, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	js_value := string(j_value)
	if f.ensure_ascii {
		js_value = escape.EscapeToASCII(js_value)
	}
	return js_value, nil
}

func (*JsonArrayFormatter) FormatItem(tableObj *table.Table, item any) ([]string, error) {
	str, err := tableObj.FormatItemBasic(item, ",")
	if err != nil {
		return nil, err
	}
	return []string{
		"[" + str + "]",
	}, nil
}

func (*JsonArrayFormatter) PrintItems(w io.Writer, tableObj *table.Table, items iface.FormattedItemList) error {
	for i := 0; i < items.Len(); i++ {
		fmt.Fprintln(w, items.Get(i)[0])
	}
	return nil
}
