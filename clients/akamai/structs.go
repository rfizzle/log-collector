package akamai

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/rfizzle/log-collector/collector"
)

type Client struct {
	collector.Client
	Options     map[string]interface{}
	domain      string
	etpConfigId string
	config      edgegrid.Config
}

type DetailBody struct {
	StartTimeSec int    `json:"startTimeSec"`
	EndTimeSec   int    `json:"endTimeSec"`
	OrderBy      string `json:"orderBy"`
	PageNumber   int    `json:"pageNumber"`
	PageSize     int    `json:"pageSize"`
	Filters      struct {
	} `json:"filters"`
}

type DetailResponse struct {
	PageInfo struct {
		TotalRecords int `json:"totalRecords"`
		PageNumber   int `json:"pageNumber"`
		PageSize     int `json:"pageSize"`
	} `json:"pageInfo"`
	DataRows []interface{} `json:"dataRows"`
}
