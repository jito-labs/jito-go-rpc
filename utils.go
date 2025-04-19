package jitorpc

import (
	"bytes"
	"encoding/json"
)

func PrettifyJSON(data json.RawMessage) string {
	var prettyJSON bytes.Buffer
	error := json.Indent(&prettyJSON, data, "", "  ")
	if error != nil {
		return string(data)
	}
	return prettyJSON.String()
}