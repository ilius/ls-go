package application

import (
	"fmt"
	"io/fs"
	"log"
	"reflect"
	"time"

	. "github.com/ilius/go-table"
	. "github.com/ilius/ls-go/common"
	"github.com/ilius/ls-go/iface"
	"github.com/ilius/ls-go/lstime"
)

var (
	t_string   = reflect.TypeOf("")
	t_uint64   = reflect.TypeOf(uint64(0))
	t_timePtr  = reflect.PtrTo(reflect.TypeOf(time.Time{}))
	t_FileMode = reflect.TypeOf(fs.FileMode(0))
)

func timeColumnFromInput(input string) string {
	switch input {
	case "", "mtime", "modified":
		return C_MTime
	case "ctime", "status", "change":
		return C_CTime
	case "atime", "access", "use", "accessed":
		return C_ATime
	case "birth", "creation", "created":
		log.Fatalf("unsupported --time=%s\n", input)
	}
	log.Fatalf("invalid --time=%s\n", input)
	return ""
}

func makeTableSpec(
	cols map[string]bool,
	formatter iface.Formatter,
	colors bool,
	nameParams *FileNameParams,
	timeParams *lstime.TimeParams,
	exprList []string,
) *TableSpec {
	tableSpec := NewTableSpec()
	if cols[C_Inode] {
		tableSpec.AddColumn(&Column{
			Name:      C_Inode,
			Title:     "inode",
			Type:      t_uint64,
			Alignment: AlignmentRight,
			Getter:    &InodeGetter{},
		})
	}
	if cols[C_Blocks] {
		tableSpec.AddColumn(&Column{
			Name:       C_Blocks,
			Title:      "Blocks",
			ShortTitle: "#B",
			Type:       t_uint64,
			Alignment:  AlignmentRight,
			Getter:     &BlocksGetter{},
		})
	}
	if cols[C_ModeOct] {
		tableSpec.AddColumn(&Column{
			Name:      C_ModeOct,
			Title:     "Oct",
			Type:      t_FileMode,
			Alignment: AlignmentRight,
			Getter:    NewOctalModeGetter(colors),
		})
	}
	if cols[C_Mode] {
		tableSpec.AddColumn(&Column{
			Name:      C_Mode,
			Title:     "Mode",
			Type:      t_FileMode,
			Alignment: nil,
			Getter:    NewModeGetter(colors),
		})
	}
	if cols[C_HardLinks] {
		tableSpec.AddColumn(&Column{
			Name:       C_HardLinks,
			Title:      "Hard Links",
			ShortTitle: "#L",
			Type:       t_uint64,
			Alignment:  AlignmentRight,
			Getter:     &HardLinksGetter{},
		})
	}
	if cols[C_Owner] {
		tableSpec.AddColumn(&Column{
			Name:       C_Owner,
			Title:      "Owner",
			ShortTitle: "O", // or "U"?
			Type:       t_string,
			Alignment:  AlignmentCenter,
			Getter:     NewOwnerGetter(colors),
		})
	}
	if cols[C_Group] {
		tableSpec.AddColumn(&Column{
			Name:       C_Group,
			Title:      "Group",
			ShortTitle: "G",
			Type:       t_string,
			Alignment:  AlignmentCenter,
			Getter:     NewGroupGetter(colors),
		})
	}
	if cols[C_Size] {
		tableSpec.AddColumn(&Column{
			Name:      C_Size,
			Title:     "Size",
			Type:      t_uint64,
			Alignment: AlignmentRight,
			Getter:    NewSizeGetter(colors, formatter.SizeFormat()),
		})
	}
	if cols[C_MTime] {
		tableSpec.AddColumn(&Column{
			Name:      C_MTime,
			Title:     "Modified Time",
			Type:      t_timePtr,
			Alignment: AlignmentRight,
			Getter:    NewMTimeGetter(colors, timeParams),
		})
	}
	if cols[C_CTime] {
		tableSpec.AddColumn(&Column{
			Name:      C_CTime,
			Title:     "Change Time",
			Type:      t_timePtr,
			Alignment: AlignmentRight,
			Getter:    NewCTimeGetter(colors, timeParams),
		})
	}
	if cols[C_ATime] {
		tableSpec.AddColumn(&Column{
			Name:      C_ATime,
			Title:     "Access Time",
			Type:      t_timePtr,
			Alignment: AlignmentRight,
			Getter:    NewATimeGetter(colors, timeParams),
		})
	}
	if cols[C_Name] {
		tableSpec.AddColumn(&Column{
			Name:      C_Name,
			Title:     "Name",
			Type:      t_string,
			Alignment: AlignmentLeft,
			Getter:    NewFileNameGetter(colors, nameParams),
		})
	}
	if len(exprList) > 0 {
		for i, exprStr := range exprList {
			getter := NewExprGetter(colors, exprStr)
			_type, err := getter.Type()
			check(err)
			alignment, err := getter.Alignment()
			check(err)
			tableSpec.AddColumn(&Column{
				Name:      fmt.Sprintf("expr%d", i+1),
				Title:     exprStr,
				Type:      _type,
				Alignment: alignment,
				Getter:    getter,
			})
		}
	}
	return tableSpec
}
