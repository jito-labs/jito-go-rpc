package jitorpc

import (
	"bytes"
	"encoding/json"
)

// PrettifyJSON formats a JSON raw message into a pretty-printed string.
func PrettifyJSON(data json.RawMessage) string {
	var prettyJSON bytes.Buffer
	error := json.Indent(&prettyJSON, data, "", "  ")
	if error != nil {
		return string(data)
	}
	return prettyJSON.String()
}
