package terminal

import (
	"fmt"
	"os"

	"github.com/ilius/consolesize-go"
)

func NewLocalTerminal() *LocalTerminal {
	return &LocalTerminal{}
}

type LocalTerminal struct {
	colorsEnabled *bool // set on first ColorsEnabled() call, used next times
}

func (*LocalTerminal) TermWidth() (int, error) {
	width, _ := consolesize.GetConsoleSize()
	return width, nil
}

func (fe *LocalTerminal) ColorsEnabled(colorFlag string) (bool, error) {
	if fe.colorsEnabled != nil {
		return *fe.colorsEnabled, nil
	}
	if os.Getenv("NO_COLOR") != "" {
		return false, nil
	}
	switch colorFlag {
	case "always", "", "y", "yes":
		return true, nil
	case "never", "n", "no":
		return false, nil
	case "auto":
		return fe.OutputIsTerminal(os.Stdout), nil
	}
	return false, fmt.Errorf("invalid --color=%s", colorFlag)
}

// check if standard output is connected to a terminal
func (f *LocalTerminal) OutputIsTerminal(stdout *os.File) bool {
	o, _ := stdout.Stat()
	if (o.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
		return true
	}
	// Pipe
	return false
}
