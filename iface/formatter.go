package iface

import (
	"io"

	"github.com/ilius/go-table"
	"github.com/ilius/ls-go/common"
	"github.com/ilius/ls-go/lscolors"
)

type FormattedItemList interface {
	Len() int
	Get(int) []string
}

type Formatter interface { //iface:ignore=unused
	// LongLinkTarget returns true if symlink targets are to be shown with --long
	LongLinkTarget() bool

	// LinkTargetSep returns separator between filename and link target with --links
	LinkTargetSep() string

	// DefaultTimeFormat returns default time format
	DefaultTimeStyle() string

	// SizeFormat returns format for file sizes
	SizeFormat() common.SizeFormat

	// PrintError formats and prints any error
	PrintError(w io.Writer, err error)

	// FolderHeader formats and prints folder header
	// when we list out any subdirectories, print those paths conspicuously
	// above the contents. this helps with visual separation
	FolderHeader(w io.Writer, path string, itemCount int)

	//  FolderHeader formats and prints folder tail
	FolderTail(w io.Writer, path string)

	// TableHeader formats and prints table header
	TableHeader(w io.Writer, tableObj *table.Table)

	// Colorize adds color to given cell/string
	Colorize(str string, style *lscolors.Style) string

	// FormatValue formats a cell value
	// used for csv and json serualization
	FormatValue(colName string, value any) (string, error)

	// FormatItem formats a file/directory item/entry
	// it does not apply table alignments
	FormatItem(tableObj *table.Table, item any) ([]string, error)

	// PrintItems applies table alignments to a list of formatted file items
	// and prints them to given io.Writer
	PrintItems(w io.Writer, tableObj *table.Table, items FormattedItemList) error
}
