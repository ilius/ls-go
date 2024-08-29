package application

import (
	"fmt"
	"os"
	"strconv"

	c "github.com/ilius/ls-go/common"
)

type SizeGetterPlain struct {
	format c.SizeFormat
}

// use metric system (or SI) to format size, powers of 1000
func (f *SizeGetterPlain) sizeStringMetric(size uint64) string {
	if size < 1000 {
		sizeStr := strconv.FormatUint(size, 10)
		return sizeStr + "B "
	}
	for _, unit := range sizeUnits {
		if unit.Next == nil || unit.Next.Metric > size {
			sizeStr := formatSizeByBase(size, unit.Metric)
			return sizeStr + unit.Symbol + " "
		}
	}
	return strconv.FormatUint(size, 10)
}

// use powers of 1024
func (f *SizeGetterPlain) sizeStringLegacy(size uint64) string {
	if size < 1024 {
		sizeStr := strconv.FormatUint(size, 10)
		return sizeStr + "B "
	}
	for _, unit := range sizeUnits {
		if unit.Next == nil || unit.Next.Legacy > size {
			sizeStr := formatSizeByBase(size, unit.Legacy)
			return sizeStr + unit.Symbol + " "
		}
	}
	return strconv.FormatUint(size, 10)
}

func (f *SizeGetterPlain) Value(item any) (any, error) {
	info, ok := item.(FileInfo)
	if !ok {
		return "", fmt.Errorf("Value: invalid type %T, must be FileInfo", item)
	}
	return uint64(info.Size()), nil
}

func (f *SizeGetterPlain) ValueString(colName string, item any) (string, error) {
	info, ok := item.(FileInfo)
	if !ok {
		return "", fmt.Errorf("ValueString: invalid type %T, must be FileInfo", item)
	}
	return app.FormatValue(colName, info.Size())
}

func (f *SizeGetterPlain) Format(item any, value any) (string, error) {
	info, ok := item.(FileInfo)
	if !ok {
		return "", fmt.Errorf("Format: invalid type %T, must be FileInfo", item)
	}
	if info.Mode()&os.ModeDevice != 0 {
		str, _ := info.DeviceNumbers()
		return str + " ", nil
	}
	size, ok := value.(uint64)
	if !ok {
		return "", fmt.Errorf("Format: invalid value type %T, must be uint64", value)
	}
	switch f.format {
	case c.SizeFormatInteger:
		return strconv.FormatUint(size, 10), nil
	case c.SizeFormatMetric:
		return f.sizeStringMetric(size), nil
	case c.SizeFormatLegacy:
		return f.sizeStringLegacy(size), nil
	}
	return "", fmt.Errorf("invalid size format %v", f.format)
}
