package application

import (
	"fmt"
	"io/fs"
	"log"
	"reflect"
	"time"

	"github.com/ilius/go-table"
	c "github.com/ilius/ls-go/common"
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
		return c.C_MTime
	case "ctime", "status", "change":
		return c.C_CTime
	case "atime", "access", "use", "accessed":
		return c.C_ATime
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
) *table.TableSpec {
	tableSpec := table.NewTableSpec()
	if cols[c.C_Inode] {
		tableSpec.AddColumn(&table.Column{
			Name:      c.C_Inode,
			Title:     "inode",
			Type:      t_uint64,
			Alignment: table.AlignmentRight,
			Getter:    &InodeGetter{},
		})
	}
	if cols[c.C_Blocks] {
		tableSpec.AddColumn(&table.Column{
			Name:       c.C_Blocks,
			Title:      "Blocks",
			ShortTitle: "#B",
			Type:       t_uint64,
			Alignment:  table.AlignmentRight,
			Getter:     &BlocksGetter{},
		})
	}
	if cols[c.C_ModeOct] {
		tableSpec.AddColumn(&table.Column{
			Name:      c.C_ModeOct,
			Title:     "Oct",
			Type:      t_FileMode,
			Alignment: table.AlignmentRight,
			Getter:    NewOctalModeGetter(colors),
		})
	}
	if cols[c.C_Mode] {
		tableSpec.AddColumn(&table.Column{
			Name:      c.C_Mode,
			Title:     "Mode",
			Type:      t_FileMode,
			Alignment: nil,
			Getter:    NewModeGetter(colors),
		})
	}
	if cols[c.C_HardLinks] {
		tableSpec.AddColumn(&table.Column{
			Name:       c.C_HardLinks,
			Title:      "Hard Links",
			ShortTitle: "#L",
			Type:       t_uint64,
			Alignment:  table.AlignmentRight,
			Getter:     &HardLinksGetter{},
		})
	}
	if cols[c.C_Owner] {
		tableSpec.AddColumn(&table.Column{
			Name:       c.C_Owner,
			Title:      "Owner",
			ShortTitle: "O", // or "U"?
			Type:       t_string,
			Alignment:  table.AlignmentCenter,
			Getter:     NewOwnerGetter(colors),
		})
	}
	if cols[c.C_Group] {
		tableSpec.AddColumn(&table.Column{
			Name:       c.C_Group,
			Title:      "Group",
			ShortTitle: "G",
			Type:       t_string,
			Alignment:  table.AlignmentCenter,
			Getter:     NewGroupGetter(colors),
		})
	}
	if cols[c.C_Size] {
		tableSpec.AddColumn(&table.Column{
			Name:      c.C_Size,
			Title:     "Size",
			Type:      t_uint64,
			Alignment: table.AlignmentRight,
			Getter:    NewSizeGetter(colors, formatter.SizeFormat()),
		})
	}
	if cols[c.C_MTime] {
		tableSpec.AddColumn(&table.Column{
			Name:      c.C_MTime,
			Title:     "Modified Time",
			Type:      t_timePtr,
			Alignment: table.AlignmentRight,
			Getter:    NewMTimeGetter(colors, timeParams),
		})
	}
	if cols[c.C_CTime] {
		tableSpec.AddColumn(&table.Column{
			Name:      c.C_CTime,
			Title:     "Change Time",
			Type:      t_timePtr,
			Alignment: table.AlignmentRight,
			Getter:    NewCTimeGetter(colors, timeParams),
		})
	}
	if cols[c.C_ATime] {
		tableSpec.AddColumn(&table.Column{
			Name:      c.C_ATime,
			Title:     "Access Time",
			Type:      t_timePtr,
			Alignment: table.AlignmentRight,
			Getter:    NewATimeGetter(colors, timeParams),
		})
	}
	if cols[c.C_Name] {
		tableSpec.AddColumn(&table.Column{
			Name:      c.C_Name,
			Title:     "Name",
			Type:      t_string,
			Alignment: table.AlignmentLeft,
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
			tableSpec.AddColumn(&table.Column{
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
