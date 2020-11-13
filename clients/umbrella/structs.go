package umbrella

import (
	"github.com/go-resty/resty/v2"
)

type AuthResponse struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type ActivityResponse struct {
	Meta interface{}   `json:"meta"`
	Data []interface{} `json:"data"`
}

type Client struct {
	Options     map[string]interface{}
	AccessToken string `json:"access_token"`
	restyClient *resty.Client
}
