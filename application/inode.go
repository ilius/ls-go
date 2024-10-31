package application

import (
	"fmt"
	"strconv"
)

type InodeGetter struct{}

func (f *InodeGetter) Value(item any) (any, error) {
	info, ok := item.(FileInfo)
	if !ok {
		return "", fmt.Errorf("Value: invalid type %T, must be FileInfo", item)
	}
	inode, err := info.Inode()
	if err != nil {
		app.AddError(err)
		return uint64(0), nil
	}
	return inode, nil
}

func (f *InodeGetter) ValueString(colName string, item any) (string, error) {
	info, ok := item.(FileInfo)
	if !ok {
		return "", fmt.Errorf("ValueString: invalid type %T, must be FileInfo", item)
	}
	inode, err := info.Inode()
	if err != nil {
		app.AddError(err)
		return "", nil
	}
	return app.FormatValue(colName, inode)
}

func (f *InodeGetter) Format(_ any, value any) (string, error) {
	// _: item is FileInfo, value is uint64 returned by .Value(item)
	return strconv.FormatUint(value.(uint64), 10), nil
}
