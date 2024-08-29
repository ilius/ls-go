package application

import (
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/ilius/go-table"
	c "github.com/ilius/ls-go/common"
)

func NewSizeGetter(colors bool, format c.SizeFormat) table.Getter {
	if colors {
		return &SizeGetter{format: format}
	}
	return &SizeGetterPlain{format: format}
}

var sizeUnits = []string{"B", "K", "M", "G", "T"}

type SizeGetter struct {
	format c.SizeFormat
}

func (f *SizeGetter) formatFloat(sizeF float64) string {
	// math.Mod(sizeF*10, 1) == 0 does not work!
	if math.Mod(sizeF*10, 1) == 0 {
		return strconv.FormatFloat(sizeF, 'f', 1, 64)
	}
	return strconv.FormatFloat(sizeF, 'f', 2, 64)
}

// use metric system (or SI) to format size, powers of 1000
func (f *SizeGetter) sizeStringMetric(size uint64) string {
	sizeFloat := float64(size)
	for i, unit := range sizeUnits {
		base := math.Pow(1000, float64(i))
		if sizeFloat < base*1000 {
			var sizeStr string
			if i == 0 {
				sizeStr = strconv.FormatUint(size, 10)
			} else {
				sizeStr = f.formatFloat(sizeFloat / base)
			}
			return app.Colorize(sizeStr+unit+" ", colors.Size[unit])
		}
	}
	return strconv.FormatUint(size, 10)
}

// use powers of 1024
func (f *SizeGetter) sizeStringLegacy(size uint64) string {
	sizeFloat := float64(size)
	for i, unit := range sizeUnits {
		base := math.Pow(1024, float64(i))
		if sizeFloat < base*1024 {
			var sizeStr string
			if i == 0 {
				sizeStr = strconv.FormatUint(size, 10)
			} else {
				sizeF := sizeFloat / base
				sizeStr = f.formatFloat(sizeF)
			}
			return app.Colorize(sizeStr+unit+" ", colors.Size[unit])
		}
	}
	return strconv.FormatUint(size, 10)
}

func (f *SizeGetter) Value(item any) (any, error) {
	info, ok := item.(FileInfo)
	if !ok {
		return "", fmt.Errorf("Value: invalid type %T, must be FileInfo", item)
	}
	return uint64(info.Size()), nil
}

func (f *SizeGetter) ValueString(colName string, item any) (string, error) {
	info, ok := item.(FileInfo)
	if !ok {
		return "", fmt.Errorf("ValueString: invalid type %T, must be FileInfo", item)
	}
	return app.FormatValue(colName, info.Size())
}

func (f *SizeGetter) Format(item any, value any) (string, error) {
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
