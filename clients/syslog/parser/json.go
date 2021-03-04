package parser

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

const (
	isoDate = `([\s]+)?(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d\.\d+([+-][0-2]\d:[0-5]\d|Z))|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d([+-][0-2]\d:[0-5]\d|Z))|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d([+-][0-2]\d:[0-5]\d|Z))`
)

// ParseJson will convert raw string received by syslog to JSON
func ParseJson(event string) ([]byte, error) {
	// Get rid of beginning date for some syslog messages
	regex := regexp.MustCompile(`^` + isoDate)
	event = regex.ReplaceAllString(event, "")

	// Get rid of beginning server name for some syslog messages
	regex = regexp.MustCompile(`^([\s]+)?[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*`)
	event = regex.ReplaceAllString(event, "")

	// Trim message
	event = strings.Trim(event, " ")

	if !isJSON(event) {
		log.Debug(event)
		return nil, fmt.Errorf("string is not in json format")
	}
	jsonMessage := []byte(event)
	return jsonMessage, nil
}

func isJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}
