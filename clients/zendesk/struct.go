package zendesk

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/rfizzle/log-collector/collector"
)

type Client struct {
	collector.Client
	Options     map[string]interface{}
	Domain      string `json:"domain"`
	restyClient *resty.Client
}

type AuditResponse struct {
	AuditLogs    []json.RawMessage `json:"audit_logs"`
  Meta      struct {
    HasMore      bool   `json:"has_more"`
    AfterCursor  string `json:"after_cursor"`
    BeforeCursor string `json:"before_cursor"`
  } `json:"meta"`
  Links struct {
    Prev string `json:"prev"`
    Next string `json:"next"`
  } `json:"links"`
}
