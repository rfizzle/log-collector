package gsuite

import (
	"github.com/rfizzle/log-collector/collector"
	"net/http"
)

type Client struct {
	collector.Client
	Options    map[string]interface{}
	httpClient *http.Client
}
