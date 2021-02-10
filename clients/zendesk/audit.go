package zendesk

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/pretty"
	"net/url"
	"strconv"
)

func (zendeskClient *Client) auditLogs(startTime, endTime string, resultsChannel chan<- string) (int, error) {
	count := 0
	page := 0
	hasNext := true

	for hasNext {
		// Increment page
		page += 1

		params := url.Values{}
		params.Add("filter[created_at][]", startTime)
		params.Add("filter[created_at][]", endTime)
		params.Set("page", strconv.Itoa(page))

		// Set URL
		scanUrl := fmt.Sprintf("https://%s/api/v2/audit_logs.json", zendeskClient.Domain)

		req := zendeskClient.restyClient.R()
		req = req.SetQueryParamsFromValues(params)
		resp, err := req.Get(scanUrl)

		if err != nil {
			return count, err
		}

		var auditResponse AuditResponse
		err = json.Unmarshal(resp.Body(), &auditResponse)

		if err != nil {
			return count, err
		}

		// Send events to results channel
		for _, event := range auditResponse.AuditLogs {
			count++
			resultsChannel <- string(pretty.Ugly(event))
		}

		hasNext = auditResponse.NextPage != ""
	}

	return count, nil
}
