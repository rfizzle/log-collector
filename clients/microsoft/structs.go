package microsoft

import (
	"encoding/json"
	"github.com/rfizzle/log-collector/collector"
	"net/http"
)

type GraphAuthResponse struct {
	TokenType    string      `json:"token_type"`
	ExpiresIn    json.Number `json:"expires_in"`
	ExtExpiresIn json.Number `json:"ext_expires_in"`
	AccessToken  string      `json:"access_token"`
}

type Client struct {
	collector.Client
	Options     map[string]interface{}
	AccessToken string `json:"access_token"`
	httpClient  *http.Client
}

type GraphSecurityAlertsResponse struct {
	Context  string        `json:"@odata.context"`
	NextLink string        `json:"@odata.nextLink"`
	Value    []interface{} `json:"value"`
}
