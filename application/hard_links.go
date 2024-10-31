package application

import (
	"fmt"
	"strconv"
)

type HardLinksGetter struct{}

func (f *HardLinksGetter) Value(item any) (any, error) {
	info, ok := item.(FileInfo)
	if !ok {
		return "", fmt.Errorf("Value: invalid type %T, must be FileInfo", item)
	}
	count, err := info.NumberOfHardLinks()
	if err != nil {
		app.AddError(err)
		return uint64(0), nil
	}
	return count, nil
}

func (f *HardLinksGetter) ValueString(colName string, item any) (string, error) {
	info, ok := item.(FileInfo)
	if !ok {
		return "", fmt.Errorf("ValueString: invalid type %T, must be FileInfo", item)
	}
	count, err := info.NumberOfHardLinks()
	if err != nil {
		app.AddError(err)
		return "", nil
	}
	return app.FormatValue(colName, count)
}

func (f *HardLinksGetter) Format(_ any, value any) (string, error) {
	// _: item is FileInfo, value is uint64 returned by .Value(item)
	return strconv.FormatUint(value.(uint64), 10), nil
}
