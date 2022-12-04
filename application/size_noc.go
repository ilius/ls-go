package application

import (
	"fmt"
	"math"
	"os"
	"strconv"

	. "github.com/ilius/ls-go/common"
)

type SizeGetterPlain struct {
	format SizeFormat
}

// use metric system (or SI) to format size, powers of 1000
func (f *SizeGetterPlain) sizeStringMetric(size uint64) string {
	sizeFloat := float64(size)
	for i, unit := range sizeUnits {
		base := math.Pow(1000, float64(i))
		if sizeFloat < base*1000 {
			var sizeStr string
			if i == 0 {
				sizeStr = strconv.FormatUint(size, 10)
			} else {
				sizeStr = strconv.FormatFloat(sizeFloat/base, 'f', 2, 64)
			}
			return sizeStr + unit + " "
		}
	}
	return strconv.FormatUint(size, 10)
}

// use powers of 1024
func (f *SizeGetterPlain) sizeStringLegacy(size uint64) string {
	sizeFloat := float64(size)
	for i, unit := range sizeUnits {
		base := math.Pow(1024, float64(i))
		if sizeFloat < base*1024 {
			var sizeStr string
			if i == 0 {
				sizeStr = strconv.FormatUint(size, 10)
			} else {
				sizeStr = strconv.FormatFloat(sizeFloat/base, 'f', 2, 64)
			}
			return sizeStr + unit + " "
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
	case SizeFormatInteger:
		return strconv.FormatUint(size, 10), nil
	case SizeFormatMetric:
		return f.sizeStringMetric(size), nil
	case SizeFormatLegacy:
		return f.sizeStringLegacy(size), nil
	}
	return "", fmt.Errorf("invalid size format %v", f.format)
}
