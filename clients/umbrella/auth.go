package umbrella

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
)

func (umbrellaClient *Client) login() error {
	// Write debug log
	log.Debugf("authenticating umbrella client...")

	// Make request
	resp, err := umbrellaClient.restyClient.R().
		SetBasicAuth(umbrellaClient.Options["key"].(string), umbrellaClient.Options["secret"].(string)).
		SetFormData(map[string]string{
			"grant_type": "client_credentials",
		}).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		Post("https://management.api.umbrella.com/auth/v2/oauth2/token")

	// Handle error
	if err != nil {
		return err
	}

	// Handle non 200
	if resp.StatusCode() != 200 {
		return fmt.Errorf("error during auth request: %s", resp.Status())
	}

	// Parse body
	var response AuthResponse
	err = json.Unmarshal(resp.Body(), &response)

	// Handle error
	if err != nil {
		return err
	}

	// Setting access token
	umbrellaClient.AccessToken = response.AccessToken

	return nil
}
