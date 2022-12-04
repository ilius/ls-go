package parse

import (
	"io/fs"
	"testing"

	"github.com/ilius/is/v2"
)

func TestParseMode(t *testing.T) {
	is := is.New(t)
	test := func(str string, mode fs.FileMode) {
		actual, err := ParseMode(str)
		if !is.NotErr(err) {
			return
		}
		is.Msg("str=%#v", str).Equal(actual, mode)
	}

	test("-rw-r--r--", 0o644)
	test("drwxr-xr-x", 0o755|fs.ModeDir)
	test("crw-------", 0o600|fs.ModeDevice|fs.ModeCharDevice)
	test("crw--w----", 0o620|fs.ModeDevice|fs.ModeCharDevice)
	test("crw-r-----", 0o640|fs.ModeDevice|fs.ModeCharDevice)
	test("crw-r--r--", 0o644|fs.ModeDevice|fs.ModeCharDevice)
	test("brw-rw----", 0o660|fs.ModeDevice)
	test("lrwxrwxrwx", 0o777|fs.ModeSymlink)

	// {"mode_oct":"4000777","mode":"drwxrwxrwt","name":"/dev/mqueue"}
	// {"mode_oct":"4000777","mode":"drwxrwxrwt","name":"/dev/shm"}
}
