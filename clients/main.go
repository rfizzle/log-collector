package clients

import (
	"errors"
	akamaiClient "github.com/rfizzle/log-collector/clients/akamai"
	gsuiteClient "github.com/rfizzle/log-collector/clients/gsuite"
	msGraph "github.com/rfizzle/log-collector/clients/microsoft"
	oktaClient "github.com/rfizzle/log-collector/clients/okta"
	"github.com/rfizzle/log-collector/collector"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
)

// Init Params

func InitClientParams() {
	flag.String("microsoft-tenant-id", "", "tenant id")
	flag.String("microsoft-client-id", "", "client id")
	flag.String("microsoft-client-secret", "", "client secret")
	flag.String("okta-domain", "", "okta domain for organization")
	flag.String("okta-api-key", "", "okta api key for authentication")
	flag.String("gsuite-credentials", "", "google service account credential file path")
	flag.String("gsuite-impersonated-user", "", "gsuite user to impersonate for API access")
	flag.String("akamai-domain", "", "akamai domain")
	flag.String("akamai-client-token", "", "akamai client token")
	flag.String("akamai-client-secret", "", "akamai client secret")
	flag.String("akamai-access-token", "", "akamai access token")
	flag.String("akamai-config-id", "", "akamai config id")
}

func InitializeClient() (collector.Client, error) {
	options, err := validateClientParams()
	if err != nil {
		return nil, err
	}
	switch viper.GetString("input") {
	case "microsoft":
		return msGraph.New(options)
	case "okta":
		return oktaClient.New(options)
	case "gsuite":
		return gsuiteClient.New(options)
	case "akamai":
		return akamaiClient.New(options)
	}
	return nil, nil
}

func validateClientParams() (map[string]interface{}, error) {
	clientOptions := make(map[string]interface{}, 0)
	switch viper.GetString("input") {
	case "microsoft":
		if viper.GetString("microsoft-tenant-id") == "" {
			return nil, errors.New("missing microsoft tenant id param (--microsoft-tenant-id)")
		}
		if viper.GetString("microsoft-client-id") == "" {
			return nil, errors.New("missing microsoft client id param (--microsoft-client-id)")
		}
		if viper.GetString("microsoft-client-secret") == "" {
			return nil, errors.New("missing microsoft client secret param (--microsoft-client-secret)")
		}
		clientOptions["tenantId"] = viper.GetString("microsoft-tenant-id")
		clientOptions["clientId"] = viper.GetString("microsoft-client-id")
		clientOptions["clientSecret"] = viper.GetString("microsoft-client-secret")
	case "okta":
		if viper.GetString("okta-domain") == "" {
			return nil, errors.New("missing okta domain param (--okta-domain)")
		}
		if viper.GetString("okta-api-key") == "" {
			return nil, errors.New("missing okta api key param (--okta-api-key)")
		}
		clientOptions["domain"] = viper.GetString("okta-domain")
		clientOptions["api-key"] = viper.GetString("okta-api-key")
	case "gsuite":
		if viper.GetString("gsuite-credentials") == "" {
			return nil, errors.New("missing google credentials param (--gsuite-credentials)")
		}
		if !fileExists(viper.GetString("gsuite-credentials")) {
			return nil, errors.New("invalid path to google credentials (--gsuite-credentials)")
		}
		if viper.GetString("gsuite-impersonated-user") == "" {
			return nil, errors.New("missing gsuite impersonate user param (--gsuite-impersonated-user)")
		}
		clientOptions["credentialFile"] = viper.GetString("gsuite-credentials")
		clientOptions["impersonationUser"] = viper.GetString("gsuite-impersonated-user")
	case "akamai":
		if viper.GetString("akamai-domain") == "" {
			return nil, errors.New("missing akamai domain param (--akamai-domain)")
		}
		if viper.GetString("akamai-client-token") == "" {
			return nil, errors.New("missing akamai client token param (--akamai-client-token)")
		}
		if viper.GetString("akamai-client-secret") == "" {
			return nil, errors.New("missing akamai client secret param (--akamai-client-secret)")
		}
		if viper.GetString("akamai-access-token") == "" {
			return nil, errors.New("missing akamai access token param (--akamai-access-token)")
		}
		if viper.GetString("akamai-config-id") == "" {
			return nil, errors.New("missing akamai config id param (--akamai-config-id)")
		}
		clientOptions["domain"] = viper.GetString("akamai-domain")
		clientOptions["clientToken"] = viper.GetString("akamai-client-token")
		clientOptions["clientSecret"] = viper.GetString("akamai-client-secret")
		clientOptions["accessToken"] = viper.GetString("akamai-access-token")
		clientOptions["configId"] = viper.GetString("akamai-config-id")
	}

	return clientOptions, nil
}

// check if file exists
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
