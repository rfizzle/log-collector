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
	NextPage     string            `json:"next_page"`
	PreviousPage string            `json:"previous_page"`
	Count        int               `json:"count"`
}
