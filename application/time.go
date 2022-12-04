package application

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ilius/go-table"
	"github.com/ilius/ls-go/lscolors"
	"github.com/ilius/ls-go/lstime"
)

func NewMTimeGetter(colors bool, params *lstime.TimeParams) table.Getter {
	if colors {
		return &MTimeGetter{&TimeGetter{params}}
	}
	return &MTimeGetterPlain{&TimeGetterPlain{params}}
}

func NewCTimeGetter(colors bool, params *lstime.TimeParams) table.Getter {
	if colors {
		return &CTimeGetter{&TimeGetter{params}}
	}
	return &CTimeGetterPlain{&TimeGetterPlain{params}}
}

func NewATimeGetter(colors bool, params *lstime.TimeParams) table.Getter {
	if colors {
		return &ATimeGetter{&TimeGetter{params}}
	}
	return &ATimeGetterPlain{&TimeGetterPlain{params}}
}

type TimeGetter struct {
	*lstime.TimeParams
}

func (f *TimeGetter) format(tm *time.Time) string {
	if tm == nil {
		return ""
	}
	if tm.IsZero() {
		return ""
	}
	if f.Relative {
		return lstime.FormatDuration((*tm).Sub(*startTime))
	}
	if f.UnixFormatStr != "" {
		return Strftime(tm, f.UnixFormatStr)
	}
	return tm.Format(f.FormatStr)
}

func (f *TimeGetter) Format(item any, value any) (string, error) {
	tm, ok := value.(*time.Time)
	if !ok {
		return "", fmt.Errorf("invalid time type %T", value)
	}
	return colorizeTimeStr(f.format(tm)) + " ", nil
}

func timeWordColor(part string) *lscolors.Style {
	if strings.IndexAny(part, "0123456789") < 0 {
		return colors.Time.Word
	}
	if strings.Index(part, ":") > 0 {
		return colors.Time.NumberColon
	}
	if strings.Index(part, "/") > 0 {
		return colors.Time.NumberSlash
	}
	num, err := strconv.ParseUint(part, 10, 64)
	if err != nil {
		return nil
	}
	if num > 1970 {
		return colors.Time.Year
	}
	return colors.Time.Number
}

func colorizeTimeStr(strPlain string) string {
	words := strings.Split(strPlain, " ")
	colored := make([]string, len(words))
	for i, word := range words {
		colored[i] = app.Colorize(word, timeWordColor(word))
	}
	return strings.Join(colored, " ")
}

type MTimeGetter struct {
	*TimeGetter
}

func (f *MTimeGetter) Value(item any) (any, error) {
	info, ok := item.(FileInfo)
	if !ok {
		return "", fmt.Errorf("Value: invalid type %T, must be FileInfo", item)
	}
	_time := info.ModTime()
	return &_time, nil
}

func (f *MTimeGetter) ValueString(colName string, item any) (string, error) {
	info, ok := item.(FileInfo)
	if !ok {
		return "", fmt.Errorf("ValueString: invalid type %T, must be FileInfo", item)
	}
	_time := info.ModTime()
	return app.FormatValue(colName, f.format(&_time))
}

type CTimeGetter struct {
	*TimeGetter
}

func (f *CTimeGetter) Value(item any) (any, error) {
	info, ok := item.(FileInfo)
	if !ok {
		return "", fmt.Errorf("Value: invalid type %T, must be FileInfo", item)
	}
	_time := info.CTime()
	return _time, nil
}

func (f *CTimeGetter) ValueString(colName string, item any) (string, error) {
	info, ok := item.(FileInfo)
	if !ok {
		return "", fmt.Errorf("ValueString: invalid type %T, must be FileInfo", item)
	}
	return app.FormatValue(colName, f.format(info.CTime()))
}

type ATimeGetter struct {
	*TimeGetter
}

func (f *ATimeGetter) Value(item any) (any, error) {
	info, ok := item.(FileInfo)
	if !ok {
		return "", fmt.Errorf("Value: invalid type %T, must be FileInfo", item)
	}
	_time := info.ATime()
	return _time, nil
}

func (f *ATimeGetter) ValueString(colName string, item any) (string, error) {
	info, ok := item.(FileInfo)
	if !ok {
		return "", fmt.Errorf("ValueString: invalid type %T, must be FileInfo", item)
	}
	return app.FormatValue(colName, f.format(info.ATime()))
}
