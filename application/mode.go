package application

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/ilius/go-table"
)

var specialPermModes = [3]fs.FileMode{
	os.ModeSticky,
	os.ModeSetgid,
	os.ModeSetuid,
}

func NewModeGetter(colors bool) table.Getter {
	if colors {
		return &ModeGetter{}
	}
	return &ModeGetterPlain{}
}

type ModeGetter struct{}

func (f *ModeGetter) Value(item any) (any, error) {
	info, ok := item.(FileInfo)
	if !ok {
		return "", fmt.Errorf("Value: invalid type %T, must be FileInfo", item)
	}
	return info.Mode(), nil
}

func (f *ModeGetter) ValueString(colName string, item any) (string, error) {
	info, ok := item.(FileInfo)
	if !ok {
		return "", fmt.Errorf("ValueString: invalid type %T, must be FileInfo", item)
	}
	return app.FormatValue(colName, formatModeNoColor(info))
}

func (f *ModeGetter) Format(item any, value any) (string, error) {
	// item is FileInfo, value is o*FileMode returned by .Value(item)
	info, ok := item.(FileInfo)
	if !ok {
		return "", fmt.Errorf("Format: invalid type %T, must be FileInfo", value)
	}

	defaultColor := colors.Perm.Other.Default()
	// info.Mode().String() does not produce the same output as `ls`, so we must build that string manually
	mode := info.Mode()
	return strings.Join([]string{
		app.Colorize(fileTypeSymbol(mode), defaultColor),
		rwxString(mode, 2, getOwnerColor(info.Owner())),
		rwxString(mode, 1, getGroupColor(info.Group())),
		rwxString(mode, 0, defaultColor),
	}, ""), nil
}
