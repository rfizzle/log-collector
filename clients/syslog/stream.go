package syslog

import (
	"encoding/json"
	"github.com/rfizzle/log-collector/clients/syslog/parser"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/pretty"
)

func (syslogClient *Client) syslogStreamParsing(streamChannel chan<- string) {
	var jsonString []byte
	var err error
	parserType := syslogClient.Options["parser"].(string)

	// Loop through channel
	for logParts := range syslogClient.logChannel {
		// Define log message
		var logMessage string

		// Check all syslog types
		if logParts["content"] == nil && logParts["message"] == nil {
			continue
		}

		// Get message from syslog struct (map key depends on format)
		if logParts["content"] != nil {
			logMessage = logParts["content"].(string)
		} else {
			logMessage = logParts["message"].(string)
		}

		// Parse content
		if parserType == "grok" {
			// Construct JSON from GROK patterns
			jsonString, err = parser.ParseGrok(logMessage, syslogClient.Options["grok-pattern"].([]string))

			// Handle errors in grok parsing
			if err != nil {
				log.Warnf("unable to marshal map to JSON: %v", err)
				continue
			}
		} else if parserType == "json" {
			// Construct JSON from raw message
			jsonString, err = parser.ParseJson(logMessage)

			// Handle errors in JSON parsing
			if err != nil {
				log.Warnf("unable to parse json message: %v", err)
				continue
			}
		} else if parserType == "kv" {
			// Construct JSON from KV string
			jsonString, err = parser.ParseKV(logMessage)

			// Handle errors in KV parsing
			if err != nil {
				log.Warnf("unable to parse kv message: %v", err)
				continue
			}
		} else if parserType == "cef" {
			// Construct JSON from CEF string
			jsonString, err = parser.ParseCef(logMessage)

			// Handle errors in CEF parsing
			if err != nil {
				log.Warnf("unable to parse cef message: %v", err)
				continue
			}
		} else if parserType == "raw" {
			jsonString, err = json.Marshal(logParts)

			// Handle errors in RAW parsing
			if err != nil {
				log.Warnf("unable to parse raw message: %v", err)
				continue
			}
		}

		if parserType != "raw" && syslogClient.Options["keep-info"].(bool) {
			// Merge message and syslog info
			finalJsonMap := make(map[string]interface{})
			err = json.Unmarshal(jsonString, &finalJsonMap)

			// Handle errors in unmarshal
			if err != nil {
				log.Warnf("unable to unmarshal json results: %v", err)
				continue
			}

			// Loop through syslog info and add to final json object
			for k, v := range logParts {
				if (k == "message" || k == "content") && !syslogClient.Options["keep-message"].(bool) {
					continue
				}
				finalJsonMap[k] = v
			}

			jsonString, err = json.Marshal(finalJsonMap)

			if err != nil {
				log.Errorf("error marshalling final json: %v", err)
				continue
			}
		}

		// Handle null parse results
		if jsonString == nil {
			log.Error("parse result for syslog message resulted in nil object")
			continue
		}

		// Write to tmp log
		streamChannel <- string(pretty.Ugly(jsonString))
	}

	return
}
