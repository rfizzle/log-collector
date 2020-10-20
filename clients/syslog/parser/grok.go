package parser

import (
	"encoding/json"
	"fmt"
	"github.com/vjeantet/grok"
)

func parseEventWithGrokPatterns(event string, grokPatterns []string) (map[string]string, error) {
	// Setup grok
	g, err := grok.NewWithConfig(&grok.Config{NamedCapturesOnly: true})

	// Handle errors
	if err != nil {
		return nil, fmt.Errorf("unable to setup grok parser: %v", err)
	}

	// Setup values map
	var values map[string]string

	// Setup grok string (this will loop through all the patterns until one works or all fail)
	for _, v := range grokPatterns {
		values, err = g.Parse(v, event)
		if err == nil {
			break
		}
	}

	// If none of the patterns worked, print error and skip to next
	if err != nil {
		return nil, fmt.Errorf("unable to parse: %v", err)
	}

	return values, nil
}

func ParseGrok(event string, grokPatterns []string) ([]byte, error) {
	values, err := parseEventWithGrokPatterns(event, grokPatterns)

	if err != nil {
		return nil, err
	}

	// Marshal map to json
	return json.Marshal(values)
}
