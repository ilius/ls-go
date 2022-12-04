package csv

import (
	"bytes"
	"encoding/csv"
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

func New(args *lsargs.Arguments) *CsvFormatter {
	// printCSVHeader()
	csvBuff := bytes.NewBuffer(nil)
	return &CsvFormatter{
		args:      args,
		csvBuff:   csvBuff,
		csvWriter: csv.NewWriter(csvBuff),
	}
}

type CsvFormatter struct {
	args      *lsargs.Arguments
	csvBuff   *bytes.Buffer
	csvWriter *csv.Writer
}

func (*CsvFormatter) LongLinkTarget() bool {
	return false
}

func (*CsvFormatter) LinkTargetSep() string {
	return ";"
}

func (*CsvFormatter) DefaultTimeStyle() string {
	timeStyle := os.Getenv("LSGO_TIME_STYLE")
	if timeStyle != "" {
		return timeStyle
	}
	return "long-iso"
}

func (*CsvFormatter) SizeFormat() SizeFormat {
	return SizeFormatInteger
}

func (*CsvFormatter) FileError(w io.Writer, err error, path string) {
	cw := csv.NewWriter(w)
	record := []string{
		path,
		"",  // Mode
		"0", // HardLinks
		"",  // Owner
		"",  // Group
		"0", // Size
		"",  // MTime

		fmt.Sprintf("ERROR: %v", err),
	}
	check(cw.Write(record))
	cw.Flush()
	check(cw.Error())
}

type ErrorWithPath interface {
	GetPath() string
}

func (f *CsvFormatter) PrintError(w io.Writer, err error) {
	path := ""
	pErr, ok := err.(ErrorWithPath)
	if ok {
		path = pErr.GetPath()
	}
	f.FileError(w, err, path)
}

func (*CsvFormatter) FolderHeader(w io.Writer, path string, itemCount int) {}

func (*CsvFormatter) FolderTail(w io.Writer, path string) {}

func (f *CsvFormatter) TableHeader(w io.Writer, tableObj *table.Table) {
	if *f.args.NoHeader {
		return
	}
	h := []string{}
	for _, col := range tableObj.Columns {
		h = append(h, col.Title)
	}
	fmt.Fprintln(w, strings.Join(h, ","))
}

func (f *CsvFormatter) Colorize(str string, style *lscolors.Style) string {
	return str
}

// previously csvString
func (f *CsvFormatter) FormatValue(colName string, value any) (string, error) {
	var valueStr string
	switch valueTyped := value.(type) {
	case string:
		valueStr = valueTyped
	case uint64:
		valueStr = strconv.FormatUint(valueTyped, 10)
	default:
		valueStr = fmt.Sprintf("%v", value)
	}
	csvWriter := f.csvWriter
	err := csvWriter.Write([]string{valueStr})
	if err != nil {
		return "", err
	}
	csvWriter.Flush()
	str := strings.TrimRight(f.csvBuff.String(), "\n")
	f.csvBuff.Reset()
	return str, nil
}

func (f *CsvFormatter) FormatItem(tableObj *table.Table, item any) ([]string, error) {
	str, err := tableObj.FormatItemBasic(item, ",")
	if err != nil {
		return nil, err
	}
	return []string{str}, nil
}

func (*CsvFormatter) PrintItems(w io.Writer, tableObj *table.Table, items iface.FormattedItemList) error {
	for i := 0; i < items.Len(); i++ {
		fmt.Fprintln(w, items.Get(i)[0])
	}
	return nil
}
