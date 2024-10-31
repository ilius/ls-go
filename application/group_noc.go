package application

import (
	"fmt"
)

type GroupGetterPlain struct{}

func (f *GroupGetterPlain) Value(item any) (any, error) {
	info, ok := item.(FileInfo)
	if !ok {
		return "", fmt.Errorf("Value: invalid type %T, must be FileInfo", item)
	}
	return info.Group(), nil
}

func (f *GroupGetterPlain) ValueString(colName string, item any) (string, error) {
	info, ok := item.(FileInfo)
	if !ok {
		return "", fmt.Errorf("ValueString: invalid type %T, must be FileInfo", item)
	}
	return app.FormatValue(colName, info.Group())
}

func (f *GroupGetterPlain) Format(_ any, value any) (string, error) {
	// _: item is FileInfo, value is string returned by .Value(item)
	return value.(string), nil
}
