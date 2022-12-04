package jsonarray

import (
	"github.com/ilius/ls-go/iface"
	"github.com/ilius/ls-go/lsargs"
)

func init() {
	var _ iface.Formatter = New(&lsargs.Arguments{})
}
