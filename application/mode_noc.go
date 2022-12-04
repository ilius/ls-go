package application

import (
	"fmt"
)

type ModeGetterPlain struct{}

func (f *ModeGetterPlain) Value(item any) (any, error) {
	info, ok := item.(FileInfo)
	if !ok {
		return "", fmt.Errorf("Value: invalid type %T, must be FileInfo", item)
	}
	return info.Mode(), nil
}

func (f *ModeGetterPlain) ValueString(colName string, item any) (string, error) {
	info, ok := item.(FileInfo)
	if !ok {
		return "", fmt.Errorf("ValueString: invalid type %T, must be FileInfo", item)
	}
	return app.FormatValue(colName, formatModeNoColor(info))
}

func (f *ModeGetterPlain) Format(item any, value any) (string, error) {
	// item is FileInfo, value is o*FileMode returned by .Value(item)
	info, ok := item.(FileInfo)
	if !ok {
		return "", fmt.Errorf("Format: invalid type %T, must be FileInfo", value)
	}
	return formatModeNoColor(info), nil
}
