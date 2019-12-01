package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

// getConfHeaders resolves headers from config overrides
func getConfHeaders(config *RequestConfig) map[string]string {
	headers := map[string]string{}
	if config != nil {
		headers = config.Headers
	}
	return headers
}

// trimMultiline - trims multiline
func trimMultiline(multi string) string {
	res := ""
	for _, line := range strings.Split(multi, "\n") {
		res += strings.Trim(line, "\t")
	}
	return res
}

// getRelativeURL out of an absolute one
func getRelativeURL(absURL string) string {
	url, _ := url.Parse(absURL)
	return url.Path
}

// getPriorEndpoint gets endpoint before the provided part ignoring case
func getPriorEndpoint(endpoint string, part string) string {
	strLen := len(strings.Split(strings.ToLower(endpoint), strings.ToLower(part))[0])
	return endpoint[:strLen]
}

// containsMetadataType checks is byte array payload contains SP OData __metadata prop
func containsMetadataType(payload []byte) bool {
	return strings.Index(fmt.Sprintf("%s", payload), `"__metadata"`) != -1
}

// patchMetadataType patches SP OData __metadata
func patchMetadataType(payload []byte, oDataType string) []byte {
	if containsMetadataType(payload) {
		return payload
	}
	metadata := map[string]interface{}{}
	if err := json.Unmarshal(payload, &metadata); err != nil {
		return payload
	}
	metadata["__metadata"] = map[string]string{"type": oDataType}
	newPayload, err := json.Marshal(metadata)
	if err != nil {
		return payload
	}
	return newPayload
}

// patchMetadataTypeCB patches SP OData __metadata with callback results
func patchMetadataTypeCB(payload []byte, resolver func() string) []byte {
	if containsMetadataType(payload) {
		return payload
	}
	return patchMetadataType(payload, resolver())
}

// parseODataItem parses OData resp taking care of OData mode
func parseODataItem(payload []byte) []byte {
	v := &struct {
		D map[string]interface{} `json:"d"`
	}{}
	if err := json.Unmarshal(payload, &v); err != nil {
		return payload
	}
	if len(v.D) == 0 {
		return payload
	}
	res, _ := json.Marshal(v.D)
	return res
}

// parseODataCollection parses OData resp taking care of OData mode
func parseODataCollection(payload []byte) [][]byte {
	v := &struct {
		D struct {
			Results []map[string]interface{} `json:"results"`
		} `json:"d"`
	}{}
	mapRes := make([]map[string]interface{}, 0)
	if err := json.Unmarshal(payload, &v); err != nil {
		if err := json.Unmarshal(payload, &mapRes); err != nil {
			return [][]byte{payload}
		}
	} else {
		mapRes = v.D.Results
	}
	res := [][]byte{}
	for _, mapItem := range mapRes {
		r, _ := json.Marshal(mapItem)
		res = append(res, []byte(r))
	}
	return res
}

// normalizeMultiLookups normalizes verbose results for multilookup
func normalizeMultiLookups(payload []byte) []byte {
	item := map[string]interface{}{}
	if err := json.Unmarshal(payload, &item); err != nil {
		return payload
	}
	for key, val := range item {
		if reflect.TypeOf(val).Kind().String() == "map" {
			results := val.(map[string]interface{})["results"]
			if results != nil {
				item[key] = results
			}
		}
	}
	normalized, err := json.Marshal(item)
	if err != nil {
		return payload
	}
	return normalized
}
