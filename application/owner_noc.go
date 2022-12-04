package application

import (
	"fmt"
)

type OwnerGetterPlain struct{}

func (f *OwnerGetterPlain) Value(item any) (any, error) {
	info, ok := item.(FileInfo)
	if !ok {
		return "", fmt.Errorf("Value: invalid type %T, must be FileInfo", item)
	}
	return info.Owner(), nil
}

func (f *OwnerGetterPlain) ValueString(colName string, item any) (string, error) {
	info, ok := item.(FileInfo)
	if !ok {
		return "", fmt.Errorf("ValueString: invalid type %T, must be FileInfo", item)
	}
	return app.FormatValue(colName, info.Owner())
}

func (f *OwnerGetterPlain) Format(item any, value any) (string, error) {
	// item is FileInfo, value is string returned by .Value(item)
	return value.(string), nil
}
