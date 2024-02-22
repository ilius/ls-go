package json

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"

	"github.com/ilius/ls-go/iface"
	// "github.com/ilius/go-table"
)

func ParseFileInfo(jsonBytes []byte) (*FakeFileInfo, error) {
	info := NewFakeFileInfo()
	err := json.Unmarshal(jsonBytes, info)
	if err != nil {
		return nil, fmt.Errorf("invalid FileInfo json %v\nerror: %w", string(jsonBytes), err)
	}
	err = info.Prepare()
	if err != nil {
		return nil, err
	}
	return info, nil
}

type ParseResult struct {
	// TableSpec *table.TableSpec
	// Columns map[string]bool
	Files []iface.FileInfo
}

func Parse(input io.Reader) (*ParseResult, error) {
	scanner := bufio.NewScanner(input)
	// cols := map[string]bool{}
	files := []iface.FileInfo{}
	// check for possible header
	if scanner.Scan() {
		line := scanner.Bytes()
		isHeader, err := lineIsHeader(line)
		if err != nil {
			return nil, err
		}
		if !isHeader {
			info, err := ParseFileInfo(line)
			if err != nil {
				return nil, err
			}
			files = append(files, info)
		}
	}
	for scanner.Scan() {
		line := scanner.Bytes()
		info, err := ParseFileInfo(line)
		if err != nil {
			return nil, err
		}
		files = append(files, info)
		// m := map[string]any{}
		// err = json.Unmarshal(line, m)
		// if err != nil {
		// 	return nil, err
		// }
		// for key, _ := range m {
		// 	cols[key] = true
		// }
	}
	result := &ParseResult{
		// Columns: cols,
		Files: files,
	}
	return result, nil
}
