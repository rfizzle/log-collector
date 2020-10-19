package gsuite

import (
	"encoding/json"
	adminreports "google.golang.org/api/admin/reports/v1"
)

func convertActivityTypeToInterface(items []*adminreports.Activity) []string {
	var data []string
	for _, val := range items {
		// Convert item to json byte array
		plain, _ := json.Marshal(val)

		// Add string to array
		data = append(data, string(plain))
	}

	return data
}
