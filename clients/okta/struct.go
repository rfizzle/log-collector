package okta

import (
	"encoding/json"
)

type OktaResponse struct {
	Uuid                  string                    `json:"uuid"`
	Published             string                    `json:"published"`
	EventType             string                    `json:"eventType"`
	Version               string                    `json:"version"`
	Severity              string                    `json:"severity"`
	LegacyEventType       string                    `json:"legacyEventType"`
	DisplayMessage        string                    `json:"displayMessage"`
	Actor                 OktaActor                 `json:"actor"`
	Client                OktaClientObj             `json:"client"`
	Outcome               OktaOutcome               `json:"outcome"`
	Target                []OktaActor               `json:"target"`
	Transaction           OktaTransaction           `json:"transaction"`
	DebugContext          OktaDebugContext          `json:"debugContext"`
	AuthenticationContext OktaAuthenticationContext `json:"authenticationContext"`
	SecurityContext       OktaSecurityContext       `json:"securityContext"`
	Request               OktaRequest               `json:"request"`
}

type OktaActor struct {
	Id          string                 `json:"id"`
	Type        string                 `json:"type"`
	AlternateId string                 `json:"alternateId"`
	DisplayName string                 `json:"displayName"`
	DetailEntry map[string]interface{} `json:"detailEntry"`
}

type OktaClientObj struct {
	UserAgent           OktaUserAgent  `json:"userAgent"`
	GeographicalContext OktaGeoContext `json:"geographicalContext"`
	Zone                string         `json:"Zone"`
	IpAddress           string         `json:"ipAddress"`
	Device              string         `json:"device"`
	Id                  string         `json:"id"`
}

type OktaUserAgent struct {
	RawUserAgent string `json:"rawUserAgent"`
	Os           string `json:"os"`
	Browser      string `json:"browser"`
}

type OktaGeoContext struct {
	Geolocation OktaGeolocation `json:"geolocation"`
	City        string          `json:"city"`
	State       string          `json:"state"`
	Country     string          `json:"country"`
	PostalCode  string          `json:"postalCode"`
}

type OktaGeolocation struct {
	Lat json.Number `json:"lat"`
	Lon json.Number `json:"lon"`
}

type OktaOutcome struct {
	Result string `json:"result"`
	Reason string `json:"reason"`
}

type OktaTransaction struct {
	Id     string                 `json:"id"`
	Type   string                 `json:"type"`
	Detail map[string]interface{} `json:"detail"`
}

type OktaDebugContext struct {
	DebugData OktaDebugData `json:"debugData"`
}

type OktaDebugData struct {
	RequestUri        string    `json:"requestUri"`
	OriginalPrincipal OktaActor `json:"originalPrincipal"`
}

type OktaAuthenticationContext struct {
	AuthenticationProvider string     `json:"authenticationProvider"`
	CredentialProvider     string     `json:"credentialProvider"`
	CredentialType         string     `json:"credentialType"`
	Issuer                 OktaIssuer `json:"issuer"`
	ExternalSessionId      string     `json:"externalSessionId"`
	Interface              string     `json:"interface"`
}

type OktaIssuer struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type OktaSecurityContext struct {
	AsNumber json.Number `json:"asNumber"`
	AsOrg    string      `json:"asOrg"`
	Isp      string      `json:"isp"`
	Domain   string      `json:"domain"`
	IsProxy  bool        `json:"isProxy"`
}

type OktaRequest struct {
	IpChain []OktaIpChain `json:"ipChain"`
}

type OktaIpChain struct {
	Ip                  string         `json:"ip"`
	GeographicalContext OktaGeoContext `json:"geographicalContext"`
	Version             string         `json:"version"`
	Source              string         `json:"source"`
}
