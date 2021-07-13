package zendesk

import (
  "encoding/json"
  "fmt"
  log "github.com/sirupsen/logrus"
  "github.com/tidwall/pretty"
  "net/url"
)

func (zendeskClient *Client) auditLogs(startTime, endTime string, resultsChannel chan<- string) (int, error) {
	count := 0
	page := 0
	hasNext := true
	after := ""

	for hasNext {
		// Increment page
		page += 1

		params := url.Values{}
		params.Set("page[size]", "100")
		params.Add("filter[created_at][]", startTime)
		params.Add("filter[created_at][]", endTime)
		if after != "" {
		  log.Info("after: %s", after)
		  params.Set("page[after]", after)
    }

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

		hasNext = auditResponse.Meta.HasMore
    after = auditResponse.Meta.AfterCursor
	}

	return count, nil
}
