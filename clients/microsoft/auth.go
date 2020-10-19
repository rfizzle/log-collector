package microsoft

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
)

// login will get a JWT with the correct grant type for collecting logs
func (microsoftClient *Client) login() error {
	params := url.Values{}
	params.Set("scope", "https://graph.microsoft.com/.default")
	params.Set("client_id", microsoftClient.Options["clientId"].(string))
	params.Set("client_secret", microsoftClient.Options["clientSecret"].(string))
	params.Set("grant_type", "client_credentials")
	body, err := microsoftClient.conductRequestRaw("POST", fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", microsoftClient.Options["tenantId"].(string)), params, "application/x-www-form-urlencoded")

	// Handle errors
	if err != nil {
		return errors.New(string(body))
	}

	// Unmarshal response json
	var authResponse GraphAuthResponse
	err = json.Unmarshal(body, &authResponse)

	// Handle error
	if err != nil {
		return errors.New(fmt.Sprintf("error on unmarshal response body: %v", err))
	}

	// Set access token
	microsoftClient.AccessToken = authResponse.AccessToken

	return nil
}
