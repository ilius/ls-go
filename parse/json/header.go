package json

import (
	"encoding/json"

	"github.com/ilius/ls-go/parse"
)

func lineIsHeader(line []byte) (bool, error) {
	infoMap := map[string]any{}
	err := json.Unmarshal(line, &infoMap)
	if err != nil {
		return false, err
	}
	for _, pair := range parse.Titles {
		value, found := infoMap[pair[0]]
		if !found {
			continue
		}
		valueStr, isStr := value.(string)
		if isStr && valueStr == pair[1] {
			return true, nil
		}
	}
	return false, nil
}
