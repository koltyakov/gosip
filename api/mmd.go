package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/csom"
)

// Taxonomy session struct
type Taxonomy struct {
	client   *HTTPClient
	config   *RequestConfig
	endpoint string
}

// NewTaxonomy - taxonomy struct constructor function
func NewTaxonomy(client *gosip.SPClient, siteURL string, config *RequestConfig) *Taxonomy {
	return &Taxonomy{
		client:   NewHTTPClient(client),
		endpoint: siteURL,
		config:   config,
	}
}

// Stores gets term stores collection object
func (taxonomy *Taxonomy) Stores() *TermStores {
	return &TermStores{
		client:   taxonomy.client,
		config:   taxonomy.config,
		endpoint: taxonomy.endpoint,
	}
}

/* Utility methods */

func appendTaxonomyProp(props []string, prop string) []string {
	for _, p := range strings.SplitN(prop, ",", -1) {
		p = strings.Trim(p, " ")
		found := false
		for _, pp := range props {
			if pp == p {
				found = true
			}
		}
		if !found {
			props = append(props, p)
		}
	}
	return props
}

func trimTaxonomyGUID(guid string) string {
	guid = strings.ToLower(guid)
	guid = strings.Replace(guid, "/guid(", "", 1)
	guid = strings.Replace(guid, ")/", "", 1)
	return guid
}

func csomResponse(httpClient *HTTPClient, siteURL string, config *RequestConfig, csomBuilder csom.Builder) (map[string]interface{}, error) {
	csomPkg, err := csomBuilder.Compile()
	if err != nil {
		return nil, err
	}

	jsomResp, err := httpClient.ProcessQuery(siteURL, bytes.NewBuffer([]byte(csomPkg)), config)
	if err != nil {
		// Retry Terms update conflicts
		if strings.Index(err.Error(), "Term update failed because of save conflict") != -1 {
			if config == nil {
				config = &RequestConfig{}
			}
			if config.Headers == nil {
				config.Headers = map[string]string{}
			}
			retryStr, ok := config.Headers["X-Gosip-Retry"]
			if !ok {
				retryStr = "0"
			}
			retry, _ := strconv.Atoi(retryStr)
			config.Headers["X-Gosip-Retry"] = strconv.Itoa(retry + 1)
			if retry+1 <= 5 {
				return csomResponse(httpClient, siteURL, config, csomBuilder)
			}
		}
		return nil, err
	}

	var jsomRespArr []interface{}
	if err := json.Unmarshal(jsomResp, &jsomRespArr); err != nil {
		return nil, err
	}

	if jsomRespArr[len(jsomRespArr)-1] == nil {
		return nil, fmt.Errorf("object not found")
	}

	res, ok := jsomRespArr[len(jsomRespArr)-1].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("can't cast CSOM response, %v+", jsomRespArr[len(jsomRespArr)-1])
	}

	return res, nil
}

func csomRespChildItems(httpClient *HTTPClient, siteURL string, config *RequestConfig, csomBuilder csom.Builder) ([]map[string]interface{}, error) {
	data, err := csomResponse(httpClient, siteURL, config, csomBuilder)
	if err != nil {
		return nil, err
	}

	items, ok := data["_Child_Items_"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("can't get child items")
	}

	var resItems []map[string]interface{}
	for _, item := range items {
		resItem, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("can't get child item")
		}
		resItems = append(resItems, resItem)
	}

	return resItems, nil
}

func csomRespChildItemsInProp(httpClient *HTTPClient, siteURL string, config *RequestConfig, csomBuilder csom.Builder, prop string) ([]map[string]interface{}, error) {
	data, err := csomResponse(httpClient, siteURL, config, csomBuilder)
	if err != nil {
		return nil, err
	}

	propData, ok := data[prop].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("can't get property data")
	}

	items, ok := propData["_Child_Items_"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("can't get child items")
	}

	var resItems []map[string]interface{}
	for _, item := range items {
		resItem, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("can't get child item")
		}
		resItems = append(resItems, resItem)
	}

	return resItems, nil
}
