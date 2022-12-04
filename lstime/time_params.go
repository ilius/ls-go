package lstime

import "fmt"

// UnixFormatStr will be empty unless set by user with --time-style=+...
// so if UnixFormatStr is set, we ignore FormatStr
// in case of UnixFormatStr, MaxWidth is set once and used afterwards
// if Relative is true, we ignore all other fields
type TimeParams struct {
	FormatStr     string
	UnixFormatStr string
	MaxWidth      int
	Relative      bool
}

func (f *TimeParams) SetTimeStyle(style string) error {
	switch style {
	case "":
		f.FormatStr = ""
		return nil
	case "full-iso", "full":
		f.FormatStr = "2006-01-02 15:04:05.999999999 Z0700"
		return nil
	case "long-iso", "long":
		f.FormatStr = "2006-01-02 15:04"
		return nil
	case "iso":
		f.FormatStr = "01-02 15:04"
		return nil
	case "locale":
		// use https://github.com/Xuanwo/go-locale
		return fmt.Errorf("unsupported time style 'locale'")
	case "relative", "rel":
		f.Relative = true
		return nil
	}
	// now it must start with + following time format
	if style[0] != '+' {
		return fmt.Errorf("invalid time style %#v", style)
	}
	f.UnixFormatStr = style[1:]
	return nil
}
