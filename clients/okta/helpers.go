package okta

import (
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// Get link from header
// resp is the response from Okta
func getResultsOffset(resp *http.Response) string {
	// Get the next link to prevent duplicates
	for _, v := range resp.Header["Link"] {
		if strings.Contains(v, "next") {
			// Build regex match
			re := regexp.MustCompile("<(.*?)>; rel=\"next\"")

			// Find the URL inside of the string
			match := re.FindStringSubmatch(v)

			// Convert to URL
			nextUrl, _ := url.Parse(match[1])

			if nextUrl.Query()["after"] == nil {
				// Get after param from URL
				return ""
			} else {
				// Get after param from URL
				return nextUrl.Query()["after"][0]
			}
		}
	}

	return ""
}

func convertLogsToString(items []interface{}) ([]string, error) {
	var data []string
	for _, val := range items {
		// Convert item to json byte array
		plain, err := json.Marshal(val)

		// Handle error
		if err != nil {
			return nil, err
		}

		// Add string to array
		data = append(data, string(plain))
	}

	return data, nil
}
