package application

import (
	"github.com/ilius/is/v2"
	"testing"
)

func Test_filevercmp(t *testing.T) {
	is := is.New(t)

	is.Equal(filevercmp("", ""), 0)
	is.Equal(filevercmp("a", "a"), 0)
	is.True(filevercmp("a", "b") < 0)
	is.True(filevercmp("b", "a") > 0)
	is.True(filevercmp("00", "01") < 0)
	is.True(filevercmp("01", "010") < 0)
	is.True(filevercmp("9", "10") < 0)
	is.True(filevercmp("0a", "0") > 0)
}
