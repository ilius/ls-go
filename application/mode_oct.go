package application

import (
	"fmt"
	"io/fs"
	"strconv"

	"github.com/ilius/go-table"
)

func NewOctalModeGetter(_ bool) table.Getter {
	// _: colors
	return &OctModeGetter{}
}

type OctModeGetter struct{}

func (f *OctModeGetter) Value(item any) (any, error) {
	info, ok := item.(FileInfo)
	if !ok {
		return "", fmt.Errorf("Value: invalid type %T, must be FileInfo", item)
	}
	return info.Mode(), nil
}

func (f *OctModeGetter) format(mode fs.FileMode) string {
	// `int64(mode) & 0777` or `int64(mode) % 0777` gives 4 octal digits
	// mode.Perm() does the same
	// although the forth digit from right seems to be always zero.
	// but I also want to show SUID/SGID/sticky
	// without type bits (link, dir, socket etc) because that makes it too long
	// and is redundent. so I came up with `mode &^ fs.ModeType`
	// which removes type bits but preserves SUID/SGID/sticky
	// example for SUID      40000777   exa: 4777
	// example for SGID      20000777   exa: 2777
	// example for sticky     4000777   exa: 1777
	// example for link    1000000777	exa: 0777
	// example for dir    20000000755	exa: 0755
	return strconv.FormatInt(int64(mode&^fs.ModeType), 8)
}

func (f *OctModeGetter) ValueString(colName string, item any) (string, error) {
	info, ok := item.(FileInfo)
	if !ok {
		return "", fmt.Errorf("ValueString: invalid type %T, must be FileInfo", item)
	}
	return app.FormatValue(colName, f.format(info.Mode()))
}

func (f *OctModeGetter) Format(item any, value any) (string, error) {
	// item is FileInfo, value is o*FileMode returned by .Value(item)
	info, ok := item.(FileInfo)
	if !ok {
		return "", fmt.Errorf("Format: invalid type %T, must be FileInfo", value)
	}
	return f.format(info.Mode()), nil
}
