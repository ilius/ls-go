package jsonarray

import (
	"github.com/ilius/ls-go/iface"
	"github.com/ilius/ls-go/lsargs"
)

func boolPtr(v bool) *bool {
	return &v
}

func init() {
	var _ iface.Formatter = New(&lsargs.Arguments{
		ASCII: boolPtr(false),
	})
}
