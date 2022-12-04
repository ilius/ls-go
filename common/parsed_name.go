package common

import "fmt"

type ParsedName struct {
	Base   string
	Ext    string
	Suffix string
}

func (p *ParsedName) String() string {
	return fmt.Sprintf(
		"{Base: %#v, Ext: %#v, Suffix: %#v}",
		p.Base,
		p.Ext,
		p.Suffix,
	)
}
