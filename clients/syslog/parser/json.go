package parser

import (
	"encoding/json"
	"fmt"
)

// ParseJson will convert raw string to JSON
func ParseJson(event string) ([]byte, error) {
	if !isJSON(event) {
		return nil, fmt.Errorf("string is not in json format")
	}
	jsonMessage := []byte(event)
	return jsonMessage, nil
}

func isJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}
