package filesystem

import (
	"github.com/ilius/ls-go/iface"
)

func init() {
	var _ iface.FileSystem = NewLocalFileSystem()
}
