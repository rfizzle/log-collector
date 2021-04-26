package clients

import (
	"errors"
	akamaiClient "github.com/rfizzle/log-collector/clients/akamai"
	fileClient "github.com/rfizzle/log-collector/clients/file"
	gsuiteClient "github.com/rfizzle/log-collector/clients/gsuite"
	msGraph "github.com/rfizzle/log-collector/clients/microsoft"
	oktaClient "github.com/rfizzle/log-collector/clients/okta"
	pubsubClient "github.com/rfizzle/log-collector/clients/pubsub"
	"github.com/rfizzle/log-collector/clients/syslog"
	umbrellaClient "github.com/rfizzle/log-collector/clients/umbrella"
	zendeskClient "github.com/rfizzle/log-collector/clients/zendesk"
	"github.com/rfizzle/log-collector/collector"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"net"
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
	flag.String("syslog-ip", "", "syslog ip address to listen on")
	flag.Int("syslog-port", 1514, "syslog port to listen on")
	flag.String("syslog-protocol", "udp", "syslog protocol to use (tcp, udp, both)")
	flag.String("syslog-parser", "raw", "syslog parser to use for syslog messages (grok, json, kv, cef, raw)")
	flag.StringArray("syslog-grok-pattern", []string{}, "syslog grok pattern to parse logs to")
	flag.Bool("syslog-keep-info", false, "syslog keep original syslog information")
	flag.Bool("syslog-keep-message", false, "syslog keep the original syslog message")
	flag.String("umbrella-key", "", "umbrella api token key")
	flag.String("umbrella-secret", "", "umbrella api token secret")
	flag.String("umbrella-org-id", "", "umbrella organization id")
	flag.String("zendesk-domain", "", "zendesk domain")
	flag.String("zendesk-email", "", "zendesk account email")
	flag.String("zendesk-password", "", "zendesk password")
	flag.String("pubsub-input-project", "", "pubsub project ID for input")
	flag.String("pubsub-input-subscription-id", "", "pubsub subscription ID for input")
	flag.String("pubsub-input-credentials", "", "pubsub credentials file for input")
	flag.String("input-file-path", "", "pattern/glob of files to read (json newline log files)")
	flag.Bool("input-file-delete", false, "delete file after successful read")
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
	case "syslog":
		return syslog.New(options)
	case "umbrella":
		return umbrellaClient.New(options)
	case "zendesk":
		return zendeskClient.New(options)
	case "pubsub":
		return pubsubClient.New(options)
	case "file":
		return fileClient.New(options)
	}
	return nil, errors.New("invalid input")
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
		clientOptions["apiKey"] = viper.GetString("okta-api-key")
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
	case "syslog":
		if !validIPAddress(viper.GetString("syslog-ip")) {
			return nil, errors.New("invalid ip param (--syslog-ip)")
		}
		if viper.GetInt("syslog-port") < 0 || viper.GetInt("port") > 65535 {
			return nil, errors.New("invalid port param (--syslog-port)")
		}

		if !contains([]string{"tcp", "udp", "both"}, viper.GetString("syslog-protocol")) {
			return nil, errors.New("invalid protocol param (--syslog-protocol)")
		}

		if !contains([]string{"grok", "json", "kv", "cef", "raw"}, viper.GetString("syslog-parser")) {
			return nil, errors.New("invalid parser param (--syslog-parser)")
		}

		if viper.GetString("syslog-parser") == "grok" && len(viper.GetStringSlice("syslog-grok-pattern")) == 0 {
			return nil, errors.New("invalid grok-pattern param (--syslog-grok-pattern)")
		}
		clientOptions["ip"] = viper.GetString("syslog-ip")
		clientOptions["port"] = viper.GetInt("syslog-port")
		clientOptions["protocol"] = viper.GetString("syslog-protocol")
		clientOptions["parser"] = viper.GetString("syslog-parser")
		clientOptions["grokPattern"] = viper.GetStringSlice("syslog-grok-pattern")
		clientOptions["keepInfo"] = viper.GetBool("syslog-keep-info")
		clientOptions["keepMessage"] = viper.GetBool("syslog-keep-message")
	case "umbrella":
		if viper.GetString("umbrella-key") == "" {
			return nil, errors.New("invalid umbrella-key param (--umbrella-key)")
		}
		if viper.GetString("umbrella-secret") == "" {
			return nil, errors.New("invalid umbrella-secret param (--umbrella-secret)")
		}
		if viper.GetString("umbrella-org-id") == "" {
			return nil, errors.New("invalid umbrella-org-id param (--umbrella-org-id)")
		}
		clientOptions["key"] = viper.GetString("umbrella-key")
		clientOptions["secret"] = viper.GetString("umbrella-secret")
		clientOptions["orgId"] = viper.GetString("umbrella-org-id")
	case "zendesk":
		if viper.GetString("zendesk-domain") == "" {
			return nil, errors.New("invalid zendesk-domain param (--zendesk-domain)")
		}
		if viper.GetString("zendesk-email") == "" {
			return nil, errors.New("invalid zendesk-email param (--zendesk-email)")
		}
		if viper.GetString("zendesk-password") == "" {
			return nil, errors.New("invalid zendesk-password param (--zendesk-password)")
		}
		clientOptions["domain"] = viper.GetString("zendesk-domain")
		clientOptions["email"] = viper.GetString("zendesk-email")
		clientOptions["password"] = viper.GetString("zendesk-password")
	case "pubsub":
		if viper.GetString("pubsub-input-project") == "" {
			return nil, errors.New("invalid pubsub project ID param (--pubsub-input-project)")
		}
		if viper.GetString("pubsub-input-subscription-id") == "" {
			return nil, errors.New("invalid pubsub subscription ID param (--pubsub-input-subscription-id)")
		}
		if !fileExists(viper.GetString("pubsub-input-credentials")) {
			return nil, errors.New("invalid pubsub credentials file (--pubsub-input-credentials)")
		}

		clientOptions["projectID"] = viper.GetString("pubsub-input-project")
		clientOptions["subscriptionID"] = viper.GetString("pubsub-input-subscription-id")
		clientOptions["credentials"] = viper.GetString("pubsub-input-credentials")
	case "file":
		if viper.GetString("input-file-path") == "" {
			return nil, errors.New("invalid input file path param (--input-file-path)")
		}

		clientOptions["path"] = viper.GetString("input-file-path")
		clientOptions["delete"] = viper.GetBool("input-file-delete")
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

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func validIPAddress(ip string) bool {
	if net.ParseIP(ip) == nil {
		return false
	} else {
		return true
	}
}
