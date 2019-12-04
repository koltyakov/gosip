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

// parseODataCollection2 parses OData resp taking care of OData mode
func parseODataCollectionPlain(payload []byte) []byte {
	v := &struct {
		D struct {
			Results []map[string]interface{} `json:"results"`
		} `json:"d"`
	}{}
	mapRes := make([]map[string]interface{}, 0)
	if err := json.Unmarshal(payload, &v); err != nil {
		if err := json.Unmarshal(payload, &mapRes); err != nil {
			return payload
		}
	} else {
		mapRes = v.D.Results
	}
	for k, m := range mapRes {
		mapRes[k] = normalizeMultiLookupsMap(m)
	}
	res, err := json.Marshal(mapRes)
	if err != nil {
		return payload
	}
	return res
}

// normalizeMultiLookups normalizes verbose results for multilookup
func normalizeMultiLookups(payload []byte) []byte {
	item := map[string]interface{}{}
	if err := json.Unmarshal(payload, &item); err != nil {
		return payload
	}
	item = normalizeMultiLookupsMap(item)
	normalized, err := json.Marshal(item)
	if err != nil {
		return payload
	}
	return normalized
}

// normalizeMultiLookupsMap normalizes verbose results for multilookup
func normalizeMultiLookupsMap(item map[string]interface{}) map[string]interface{} {
	for key, val := range item {
		if val != nil && reflect.TypeOf(val).Kind().String() == "map" {
			results := val.(map[string]interface{})["results"]
			if results != nil {
				item[key] = results
			}
		}
	}
	return item
}

// fixDateInResponse fixes incorrect date responses for a provided fields
func fixDatesInResponse(data []byte, dateFields []string) []byte {
	metadata := map[string]interface{}{}
	if err := json.Unmarshal(data, &metadata); err != nil {
		return data
	}
	for _, k := range dateFields {
		val := metadata[k]
		if val != nil && reflect.TypeOf(val).String() == "string" {
			if len(fmt.Sprintf("%s", val)) == len("2019-12-03T12:19:45") {
				metadata[k] = fmt.Sprintf("%sZ", val)
			}
		}
	}
	res, _ := json.Marshal(metadata)
	return res
}
