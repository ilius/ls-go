package application

import (
	"fmt"
	"time"

	"github.com/ilius/ls-go/lstime"
)

type TimeGetterPlain struct {
	*lstime.TimeParams
}

func (f *TimeGetterPlain) format(tm *time.Time) string {
	if tm == nil {
		return ""
	}
	if tm.IsZero() {
		return ""
	}
	if f.Relative {
		return lstime.FormatDuration(time.Until(*tm))
	}
	if f.UnixFormatStr != "" {
		return Strftime(tm, f.UnixFormatStr)
	}
	return tm.Format(f.FormatStr)
}

func (f *TimeGetterPlain) Format(item any, value any) (string, error) {
	tm, ok := value.(*time.Time)
	if !ok {
		return "", fmt.Errorf("invalid time type %T", value)
	}
	return f.format(tm) + " ", nil
}

type MTimeGetterPlain struct {
	*TimeGetterPlain
}

func (f *MTimeGetterPlain) Value(item any) (any, error) {
	info, ok := item.(FileInfo)
	if !ok {
		return "", fmt.Errorf("Value: invalid type %T, must be FileInfo", item)
	}
	_time := info.ModTime()
	return &_time, nil
}

func (f *MTimeGetterPlain) ValueString(colName string, item any) (string, error) {
	info, ok := item.(FileInfo)
	if !ok {
		return "", fmt.Errorf("ValueString: invalid type %T, must be FileInfo", item)
	}
	_time := info.ModTime()
	return app.FormatValue(colName, f.format(&_time))
}

type CTimeGetterPlain struct {
	*TimeGetterPlain
}

func (f *CTimeGetterPlain) Value(item any) (any, error) {
	info, ok := item.(FileInfo)
	if !ok {
		return "", fmt.Errorf("Value: invalid type %T, must be FileInfo", item)
	}
	_time := info.CTime()
	return _time, nil
}

func (f *CTimeGetterPlain) ValueString(colName string, item any) (string, error) {
	info, ok := item.(FileInfo)
	if !ok {
		return "", fmt.Errorf("ValueString: invalid type %T, must be FileInfo", item)
	}
	return app.FormatValue(colName, f.format(info.CTime()))
}

type ATimeGetterPlain struct {
	*TimeGetterPlain
}

func (f *ATimeGetterPlain) Value(item any) (any, error) {
	info, ok := item.(FileInfo)
	if !ok {
		return "", fmt.Errorf("Value: invalid type %T, must be FileInfo", item)
	}
	_time := info.ATime()
	return _time, nil
}

func (f *ATimeGetterPlain) ValueString(colName string, item any) (string, error) {
	info, ok := item.(FileInfo)
	if !ok {
		return "", fmt.Errorf("ValueString: invalid type %T, must be FileInfo", item)
	}
	return app.FormatValue(colName, f.format(info.ATime()))
}
