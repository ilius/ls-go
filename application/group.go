package application

import (
	"fmt"

	"github.com/ilius/go-table"
	"github.com/ilius/ls-go/lscolors"
)

func NewGroupGetter(colors bool) table.Getter {
	if colors {
		return &GroupGetter{}
	}
	return &GroupGetterPlain{}
}

func getGroupColor(group string) *lscolors.Style {
	return colors.Perm.Group.Get(group)
}

type GroupGetter struct{}

func (f *GroupGetter) Value(item any) (any, error) {
	info, ok := item.(FileInfo)
	if !ok {
		return "", fmt.Errorf("Value: invalid type %T, must be FileInfo", item)
	}
	return info.Group(), nil
}

func (f *GroupGetter) ValueString(colName string, item any) (string, error) {
	info, ok := item.(FileInfo)
	if !ok {
		return "", fmt.Errorf("ValueString: invalid type %T, must be FileInfo", item)
	}
	return app.FormatValue(colName, info.Group())
}

func (f *GroupGetter) Format(item any, value any) (string, error) {
	// item is FileInfo, value is string returned by .Value(item)
	group := value.(string)
	return app.Colorize(group, getGroupColor(group)), nil
}
