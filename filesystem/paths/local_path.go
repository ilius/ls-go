package paths

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ilius/ls-go/common"
)

// this must not be used outside this file
const localPathSep = string(filepath.Separator)

var (
	suffixRE = regexp.MustCompile(`(\.[0-9]+)*` + "$")
	extRE    = regexp.MustCompile(`\.[a-zA-Z0-9]+$`)
)

// abstraction of "path/filepath" methods
type LocalFilePath struct{}

// same as filepath.Abs
func (*LocalFilePath) Abs(path string) (string, error) {
	return filepath.Abs(path)
}

// same as filepath.Dir, but also work with path ending with a slash/sep
func (*LocalFilePath) Dir(path string) string {
	return filepath.Dir(strings.TrimRight(path, localPathSep))
}

// checks if path is absolute
func (*LocalFilePath) IsAbs(path string) bool {
	// On Windows OS, IsAbs validates if a path is valid based on if stars with a
	// unit (eg.: `C:\`)  to assert that is absolute, but in this mem implementation
	// any path starting by `separator` is also considered absolute.
	return filepath.IsAbs(path) || strings.HasPrefix(path, localPathSep)
}

// same as filepath.Join
func (*LocalFilePath) Join(elem ...string) string {
	return filepath.Join(elem...)
}

// same as filepath.Rel
func (*LocalFilePath) Rel(basepath, targpath string) (string, error) {
	return filepath.Rel(basepath, targpath)
}

// Clean returns the shortest path name equivalent to path by purely lexical
// processing. See "path/filepath.Clean"
func (*LocalFilePath) Clean(path string) string {
	return filepath.Clean(path)
}

// Base returns the last element of path.
// Trailing path separators are removed before extracting the last element.
// If the path is empty, Base returns ".".
// If the path consists entirely of separators, Base returns a single separator.
func (*LocalFilePath) Base(path string) string {
	return filepath.Base(path)
}

// Split splits path immediately following the final Separator, separating it
// into a directory and file name component. If there is no Separator in path,
// Split returns an empty dir and file set to path. The returned values have
// the property that path = dir+file.
func (*LocalFilePath) Split(path string) (string, string) {
	return filepath.Split(path)
}

// FromSlash returns the result of replacing each slash ('/') character in
// path with a separator character. Multiple slashes are replaced by multiple separators.
func (*LocalFilePath) FromSlash(path string) string {
	return filepath.FromSlash(path)
}

func (*LocalFilePath) SplitAll(path string) []string {
	return strings.Split(path, localPathSep)
}

// same as Join, but give a color to path sep (slash or backslash)
func (*LocalFilePath) JoinColor(color string, reset string, elem ...string) string {
	return strings.Join(elem, color+localPathSep+reset)
}

// SplitExt returns {basename, ext, suffix}
// for example:
// filename="test.ZIP" 		-> {"test", ".ZIP", ""}
// filename="test.ZIP.1" 	-> {"test", ".ZIP", ".1"}
// filename="lib.so.1.2.3"	-> {"lib", ".so.1", ".1.2.3"}
func (*LocalFilePath) SplitExt(filename string) *common.ParsedName {
	suffix := suffixRE.FindString(filename)
	filename = filename[:len(filename)-len(suffix)]
	ext := extRE.FindString(filename)
	if len(ext) == len(filename) {
		ext = ""
	}
	basename := filename[:len(filename)-len(ext)]
	if strings.HasSuffix(basename, localPathSep) {
		basename = filename
		ext = ""
	}
	return &common.ParsedName{
		Base:   basename,
		Ext:    ext,
		Suffix: suffix,
	}
}
