package paths

import (
	"path/filepath"
	"testing"

	"github.com/ilius/is/v2"
)

func TestSplitExt(t *testing.T) {
	is := is.New(t)
	local := &LocalFilePath{}

	test := func(filename string, basename string, ext string, suffix string) {
		actual := local.SplitExt(filename)
		is := is.AddMsg(
			"filename=%#v, actual: %#v",
			filename, actual,
		)
		is.Equal(actual.Base, basename)
		is.Equal(actual.Ext, ext)
		is.Equal(actual.Suffix, suffix)
	}
	sep := string(filepath.Separator)

	test("test", "test", "", "")
	test("README.md", "README", ".md", "")
	test(".test", ".test", "", "")
	test(".test.txt", ".test", ".txt", "")
	test("a.test", "a", ".test", "")
	test("test.EXE", "test", ".EXE", "")
	test("test.abc-d", "test.abc-d", "", "")
	test("test.7z", "test", ".7z", "")
	test(".git", ".git", "", "")
	test("a"+sep+".git", "a"+sep+".git", "", "")
	test("."+sep+".git", "."+sep+".git", "", "")

	test("test.so", "test", ".so", "")
	test("test.so.0", "test", ".so", ".0")
	test("test.so.1", "test", ".so", ".1")
	test("test.so.2", "test", ".so", ".2")
	test("test.so.1.2.3", "test", ".so", ".1.2.3")
	test("libarmadillo.so.11", "libarmadillo", ".so", ".11")

	test("test.7z", "test", ".7z", "")
	test("test.7z.0", "test", ".7z", ".0")
	test("test.7z.1", "test", ".7z", ".1")
	test("test.7z.2", "test", ".7z", ".2")
	test("test.7z.0998", "test", ".7z", ".0998")
}
