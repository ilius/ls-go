package common

import "fmt"

type FileError struct {
	Path string `json:"name"`
	Msg  string `json:"error"`
}

func (e *FileError) Error() string {
	return fmt.Sprintf("%v, file %#v", e.Msg, e.Path)
}

func (e *FileError) GetPath() string {
	return e.Path
}
