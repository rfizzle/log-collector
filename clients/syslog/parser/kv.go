package parser

import (
	"encoding/json"
	"fmt"
	"github.com/jjeffery/kv"
)

// parseKeyValue will take a key value formatted string and convert it into a key value map
func parseKeyValue(event string, cef bool) (map[string]string, error) {
	// Use KeyValue library to parse
	text, list := kv.Parse([]byte(event))

	// If return text then an error occurred during parsing
	if len(text) > 0 {
		return nil, fmt.Errorf(`invalid key value format at: "%s"`, string(text))
	}

	// Convert from list to a map
	elementMap := make(map[string]string)
	for i := 0; i < len(list); i += 2 {
		key := list[i].(string)
		value := list[i+1].(string)
		if cef {
			elementMap[cefEscapeExtension(key)] = cefEscapeExtension(value)
		} else {
			elementMap[key] = value
		}
	}

	return elementMap, nil
}

// ConstructKeyValue will take a key value formatted string and convert it into a key value json object
func ParseKV(event string) ([]byte, error) {
	// Parse key value string
	result, err := parseKeyValue(event, false)

	// Handle errors
	if err != nil {
		return nil, err
	}

	// Marshal JSON string
	jsonString, err := json.Marshal(result)

	// Handle errors
	if err != nil {
		return nil, err
	}

	return jsonString, nil
}
