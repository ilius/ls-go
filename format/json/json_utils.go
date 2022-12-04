package json

import (
	"encoding/json"

	"github.com/ilius/ls-go/escape"
)

type ItemErrorJSON struct {
	Name  string `json:"name"`
	Error string `json:"error"`
}

func jsonKeyValue(key string, value any, ensure_ascii bool) (string, error) {
	j_key, err := json.Marshal(key)
	if err != nil {
		return "", err
	}
	j_value, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	js_value := string(j_value)
	if ensure_ascii {
		js_value = escape.EscapeToASCII(js_value)
	}
	return string(j_key) + ":" + js_value, nil
}
