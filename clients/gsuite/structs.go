package gsuite

import "net/http"

type Client struct {
	Options    map[string]interface{}
	httpClient *http.Client
}
