package application

import (
	"fmt"
	"strconv"
)

type BlocksGetter struct{}

func (f *BlocksGetter) Value(item any) (any, error) {
	info, ok := item.(FileInfo)
	if !ok {
		return "", fmt.Errorf("Value: invalid type %T, must be FileInfo", item)
	}
	return uint64(info.Blocks()), nil
}

func (f *BlocksGetter) ValueString(colName string, item any) (string, error) {
	info, ok := item.(FileInfo)
	if !ok {
		return "", fmt.Errorf("ValueString: invalid type %T, must be FileInfo", item)
	}
	return app.FormatValue(colName, info.Blocks())
}

func (f *BlocksGetter) Format(item any, value any) (string, error) {
	info, ok := item.(FileInfo)
	if !ok {
		return "", fmt.Errorf("Format: invalid type %T, must be FileInfo", item)
	}
	return strconv.FormatInt(
		info.Blocks(),
		10,
	), nil
}
