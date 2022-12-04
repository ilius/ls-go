package terminal

import "github.com/ilius/ls-go/iface"

func init() {
	var _ iface.Terminal = NewLocalTerminal()
}
