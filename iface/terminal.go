package iface

import "os"

type Terminal interface {
	// TermWidth returns terminal width
	TermWidth() (int, error)

	// OutputIsTerminal returns true if standard output is connected to a terminal
	OutputIsTerminal(stdout *os.File) bool

	// ColorsEnabled returns true if colors are enabled
	ColorsEnabled(colorFlag string) (bool, error)
}
