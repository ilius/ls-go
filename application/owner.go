package application

import (
	"fmt"

	"github.com/ilius/go-table"
	"github.com/ilius/ls-go/lscolors"
)

func NewOwnerGetter(colors bool) table.Getter {
	if colors {
		return &OwnerGetter{}
	}
	return &OwnerGetterPlain{}
}

func getOwnerColor(owner string) *lscolors.Style {
	if owner == app.Platform.UserName() {
		owner = lscolors.SELF
	}
	return colors.Perm.User.Get(owner)
}

type OwnerGetter struct{}

func (f *OwnerGetter) Value(item any) (any, error) {
	info, ok := item.(FileInfo)
	if !ok {
		return "", fmt.Errorf("Value: invalid type %T, must be FileInfo", item)
	}
	return info.Owner(), nil
}

func (f *OwnerGetter) ValueString(colName string, item any) (string, error) {
	info, ok := item.(FileInfo)
	if !ok {
		return "", fmt.Errorf("ValueString: invalid type %T, must be FileInfo", item)
	}
	return app.FormatValue(colName, info.Owner())
}

func (f *OwnerGetter) Format(_ any, value any) (string, error) {
	// _: item is FileInfo, value is string returned by .Value(item)
	owner := value.(string)
	return app.Colorize(owner, getOwnerColor(owner)), nil
}
