package tabular

import (
	"github.com/ilius/ls-go/iface"
)

func init() {
	var _ iface.Formatter = &TabularFormatter{}
}
