package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
)

func (esClient *Client) es6Search(ctx context.Context) ([]string, error) {
	var r map[string]interface{}
	// Perform the search request.
	res, err := esClient.es6Client.Search(
		esClient.es6Client.Search.WithContext(ctx),
		esClient.es6Client.Search.WithIndex(esClient.index),
		esClient.es6Client.Search.WithBody(bytes.NewReader(esClient.query)),
		esClient.es6Client.Search.WithTrackTotalHits(true),
		esClient.es6Client.Search.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Errorf("error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, fmt.Errorf("error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			return nil, fmt.Errorf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("error parsing the response body: %s", err)
	}

	return searchResults(r)
}

func (esClient *Client) es7Search(ctx context.Context) ([]string, error) {
	var r map[string]interface{}
	// Perform the search request.
	res, err := esClient.es7Client.Search(
		esClient.es7Client.Search.WithContext(ctx),
		esClient.es7Client.Search.WithIndex(esClient.index),
		esClient.es7Client.Search.WithBody(bytes.NewReader(esClient.query)),
		esClient.es7Client.Search.WithTrackTotalHits(true),
		esClient.es7Client.Search.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Errorf("error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, fmt.Errorf("error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			return nil, fmt.Errorf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("error parsing the response body: %s", err)
	}

	return searchResults(r)
}

func (esClient *Client) es8Search(ctx context.Context) ([]string, error) {
	var r map[string]interface{}
	// Perform the search request.
	res, err := esClient.es8Client.Search(
		esClient.es8Client.Search.WithContext(ctx),
		esClient.es8Client.Search.WithIndex(esClient.index),
		esClient.es8Client.Search.WithBody(bytes.NewReader(esClient.query)),
		esClient.es8Client.Search.WithTrackTotalHits(true),
		esClient.es8Client.Search.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Errorf("error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, fmt.Errorf("error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			return nil, fmt.Errorf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("error parsing the response body: %s", err)
	}

	return searchResults(r)
}

func searchResults(results map[string]interface{}) ([]string, error) {
	var arr []string
	if val, ok := results["hits"].(map[string]interface{}); ok {
		if val2, ok2 := val["hits"].([]interface{}); ok2 {
			for _, hit := range val2 {
				item := hit.(map[string]interface{})
				resBytes, err := json.Marshal(item)
				if err != nil {
					log.Errorf("issue parsing result: %v", err)
					continue
				}
				arr = append(arr, string(resBytes))
			}
		} else {
			return nil, fmt.Errorf(`"issue parsing results["hits"]["hits"]`)
		}
	} else {
		return nil, fmt.Errorf(`"issue parsing results["hits"]`)
	}

	return arr, nil
}
