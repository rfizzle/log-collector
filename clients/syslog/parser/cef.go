package parser

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/dlclark/regexp2"
)

type CefEvent struct {
	Version            string
	DeviceVendor       string
	DeviceProduct      string
	DeviceVersion      string
	DeviceEventClassId string
	Name               string
	Severity           string
	Extensions         map[string]string
}

func ParseCef(event string) ([]byte, error) {
	cefEvent, err := cefStringToObject(event)

	if err != nil {
		return nil, err
	}

	// Marshal JSON string
	jsonString, err := json.Marshal(cefEvent)

	// Handle errors
	if err != nil {
		fmt.Println(3)
		return nil, err
	}

	return jsonString, nil
}

func cefStringToObject(cefString string) (*CefEvent, error) {
	// Split by CEF separator
	arr := strings.Split(cefString, "|")

	if len(arr) < 8 {
		return nil, fmt.Errorf("invalid CEF format")
	}

	version := ""

	if strings.Contains(arr[0], ":") {
		// Split first field to validate CEF
		validate := strings.Split(arr[0], ":")

		// Validate that it is a valid CEF message
		if validate[0] != "CEF" {
			return nil, fmt.Errorf("invalid CEF format")
		}

		version = validate[1]
	} else {
		if _, err := strconv.Atoi(arr[0]); err != nil {
			return nil, fmt.Errorf("invalid CEF format")
		}
		version = arr[0]
	}

	// Get extensions
	extensions := strings.Join(arr[7:], "|")

	// Replace colons with {{COLON}}
	safeExtensions := strings.ReplaceAll(extensions, ":", "{{COLON}}")
	safeExtensions = strings.ReplaceAll(safeExtensions, `\\=`, "{{EQUAL_ESCAPE_2}}")
	safeExtensions = strings.ReplaceAll(safeExtensions, `\=`, "{{EQUAL_ESCAPE_1}}")

	// Replace non KV spaces with {{SPACE}}
	re := regexp2.MustCompile(`\s(?!([\w\-]+)\=)`, 0)
	safeExtensions2, err := re.Replace(safeExtensions, "{{SPACE}}", -1, -1)

	// Parse extensions in key value format
	keyValueMap, err := parseKeyValue(safeExtensions2, true)

	// Handle error
	if err != nil {
		return nil, err
	}

	// Restore colons
	newKeyValueMap := make(map[string]string, 0)
	for k, v := range keyValueMap {
		newKey := strings.ReplaceAll(k, `{{SPACE}}`, " ")
		newKey = strings.ReplaceAll(newKey, `{{EQUAL_ESCAPE_1}}`, `\=`)
		newKey = strings.ReplaceAll(newKey, `{{EQUAL_ESCAPE_2}}`, `\\=`)
		newKey = strings.ReplaceAll(newKey, `{{COLON}}`, ":")

		newValue := strings.ReplaceAll(v, `{{SPACE}}`, " ")
		newValue = strings.ReplaceAll(newValue, `{{EQUAL_ESCAPE_1}}`, `\=`)
		newValue = strings.ReplaceAll(newValue, `{{EQUAL_ESCAPE_2}}`, `\\=`)
		newValue = strings.ReplaceAll(newValue, `{{COLON}}`, ":")
		newValue = strings.TrimSpace(newValue)

		newKeyValueMap[newKey] = newValue
	}

	// Build CEF event
	cefEvent := &CefEvent{
		Version:            version,
		DeviceVendor:       cefEscapeField(arr[1]),
		DeviceProduct:      cefEscapeField(arr[2]),
		DeviceVersion:      cefEscapeField(arr[3]),
		DeviceEventClassId: cefEscapeField(arr[4]),
		Name:               cefEscapeField(arr[5]),
		Severity:           cefEscapeField(arr[6]),
		Extensions:         newKeyValueMap,
	}

	return cefEvent, nil
}

// Unescape CEF fields
func cefEscapeField(field string) string {

	replacer := strings.NewReplacer(
		"\\\\", "\\",
		"\\|", "|",
		"\\n", "\n",
	)

	return replacer.Replace(field)
}

// Unescape CEF extensions
func cefEscapeExtension(field string) string {

	replacer := strings.NewReplacer(
		"\\\\", "\\",
		"\\n", "\n",
		"\\=", "=",
	)

	return replacer.Replace(field)
}
