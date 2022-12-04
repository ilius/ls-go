package lsplatform

import (
	"fmt"
	"io/fs"
)

type LocalPlatform struct{}

func New() *LocalPlatform {
	return &LocalPlatform{}
}

type FileInfo interface {
	fs.FileInfo
	PathAbs() string
}

type OwnerGroup struct {
	Owner string
	Group string
}

type PlatformError struct {
	Operation string `json:"operation"`
	Path      string `json:"name"`
	Msg       string `json:"error"`
}

func (e *PlatformError) Error() string {
	return fmt.Sprintf("error in %v: %v: file %#v", e.Operation, e.Msg, e.Path)
}

func (e *PlatformError) GetPath() string {
	return e.Path
}
