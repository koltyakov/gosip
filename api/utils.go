package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"reflect"
	"strings"
)

// TrimMultiline trims multiline string converting to a single line friendly for payloads
func TrimMultiline(multi string) string {
	res := ""
	for _, line := range strings.Split(multi, "\n") {
		res += strings.Trim(line, "\t")
	}
	return res
}

// NormalizeODataItem parses OData resp taking care of OData mode
func NormalizeODataItem(payload []byte) []byte {
	v := &struct {
		D map[string]interface{} `json:"d"`
	}{}
	if err := json.Unmarshal(payload, &v); err != nil {
		return payload
	}
	if len(v.D) == 0 {
		return payload
	}
	v.D = normalizeMultiLookupsMap(v.D)
	res, _ := json.Marshal(v.D)
	return res
}

// NormalizeODataCollection parses OData resp taking care of OData mode
func NormalizeODataCollection(payload []byte) ([]byte, string) {
	bb, nextURL := normalizeODataCollection(payload)
	var mapRes []map[string]interface{}
	for _, b := range bb {
		mapItem := map[string]interface{}{}
		if err := json.Unmarshal(b, &mapItem); err == nil {
			mapRes = append(mapRes, normalizeMultiLookupsMap(mapItem))
		}
	}
	res, err := json.Marshal(mapRes)
	if err != nil {
		return payload, nextURL
	}
	return res, nextURL
}

// ExtractEntityURI extracts REST entity URI from payload
func ExtractEntityURI(payload []byte) string {
	payload = NormalizeODataItem(payload)
	r := &struct {
		Metadata *struct {
			ID string `json:"id"`
		} `json:"__metadata"`
		ID string `json:"odata.id"`
	}{}
	entityURI := ""
	if err := json.Unmarshal(payload, &r); err == nil {
		entityURI = r.ID
	}
	if r.Metadata != nil && r.Metadata.ID != "" {
		entityURI = r.Metadata.ID
	}
	return entityURI
}

// EscapePathURI escapes URI path for GetFile and GetFolder methods
func EscapePathURI(s string) string {
	s = strings.Replace(s, "'", "''", -1)
	s = strings.Replace(s, "%", "%25", -1)
	s = strings.Replace(s, "#", "%23", -1)
	return s
}

// getConfHeaders resolves headers from config overrides
func getConfHeaders(config *RequestConfig) map[string]string {
	headers := map[string]string{}
	if config != nil {
		headers = config.Headers
	}
	return headers
}

// patchConfigHeaders patches config headers and returns a copy of config
func patchConfigHeaders(config *RequestConfig, headers map[string]string) *RequestConfig {
	conf := &RequestConfig{}
	if config != nil {
		conf.Context = config.Context
		conf.Headers = config.Headers
	}
	if conf.Headers == nil {
		conf.Headers = map[string]string{}
	}
	for k, v := range headers {
		conf.Headers[k] = v
	}
	return conf
}

// getRelativeURL out of an absolute one
func getRelativeURL(absURL string) string {
	u, _ := url.Parse(absURL)
	return u.Path
}

// checkGetRelativeURL checks if URL is relative, prepends relative part if missed
func checkGetRelativeURL(relativeURI string, ctxURL string) string {
	absoluteURL := getPriorEndpoint(ctxURL, "/_api")
	relativeURL := getRelativeURL(absoluteURL)

	// Empty URI should end up with site relative URL
	if relativeURI == "" {
		return relativeURL
	}

	// A URI like "Lists/List" should be transformed to "/sites/site/Lists/List"
	if string([]rune(relativeURI)[0]) != "/" {
		relativeURI = fmt.Sprintf("%s/%s", relativeURL, relativeURI)
	}

	return relativeURI
}

// getPriorEndpoint gets endpoint before the provided part ignoring case
func getPriorEndpoint(endpoint string, part string) string {
	strLen := len(strings.Split(strings.ToLower(endpoint), strings.ToLower(part))[0])
	return endpoint[:strLen]
}

// getIncludeEndpoint gets endpoint including the provided part ignoring case
func getIncludeEndpoint(endpoint string, part string) string {
	strLen := len(strings.Split(strings.ToLower(endpoint), strings.ToLower(part))[0])
	if len(endpoint) == strLen {
		return endpoint
	}
	return endpoint[:strLen] + part
}

// getIncludeEndpoints gets endpoint including the provided parts as array of values ignoring case
func getIncludeEndpoints(endpoint string, parts []string) string {
	result := endpoint
	for _, part := range parts {
		res := getIncludeEndpoint(endpoint, part)
		if len(res) < len(endpoint) {
			result = res
		}
	}
	return result
}

// containsMetadataType checks is byte array payload contains SP OData __metadata prop
func containsMetadataType(payload []byte) bool {
	return strings.Contains(string(payload), `"__metadata"`)
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

// normalizeODataCollection parses OData resp taking care of OData mode
func normalizeODataCollection(payload []byte) ([][]byte, string) {
	r := &struct {
		// Verbose OData structure
		D struct {
			Results []map[string]interface{} `json:"results"`
			NextURL string                   `json:"__next"`
		} `json:"d"`
		// Minimalmatadata/Nometadata OData structure
		Results []map[string]interface{} `json:"value"`
		NextURL string                   `json:"odata.nextLink"`
	}{}
	mapRes := make([]map[string]interface{}, 0)
	nextURL := ""
	if err := json.Unmarshal(payload, &r); err == nil {
		mapRes = r.Results
		nextURL = r.NextURL
	}
	if r.Results == nil {
		mapRes = r.D.Results
		nextURL = r.D.NextURL
	}
	if r.D.Results == nil && r.Results == nil {
		if err := json.Unmarshal(payload, &mapRes); err != nil {
			return [][]byte{payload}, ""
		}
	}
	var res [][]byte
	for _, mapItem := range mapRes {
		r, _ := json.Marshal(mapItem)
		res = append(res, r)
	}
	return res, nextURL
}

// getODataCollectionNextPageURL parses OData resp taking care of OData mode
func getODataCollectionNextPageURL(payload []byte) string {
	r := &struct {
		// Verbose OData structure
		D struct {
			NextURL string `json:"__next"`
		} `json:"d"`
		// Minimalmatadata/Nometadata OData structure
		NextURL string `json:"odata.nextLink"`
	}{}
	if err := json.Unmarshal(payload, &r); err != nil {
		return ""
	}
	if r.NextURL != "" {
		return r.NextURL
	}
	if r.D.NextURL != "" {
		return r.D.NextURL
	}
	return ""
}

// normalizeMultiLookups normalizes verbose results for multi-lookup
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

// normalizeMultiLookupsMap normalizes verbose results for multi-lookup
func normalizeMultiLookupsMap(item map[string]interface{}) map[string]interface{} {
	for key, val := range item {
		if val != nil && reflect.TypeOf(val).Kind().String() == "map" {
			results := val.(map[string]interface{})["results"]
			if results != nil {
				item[key] = results
			}
			// Recursive props normalization
			if item[key] != nil && reflect.TypeOf(item[key]).Kind().String() == "map" {
				item[key] = normalizeMultiLookupsMap(item[key].(map[string]interface{}))
			}
			if item[key] != nil && reflect.TypeOf(item[key]).Kind().String() == "slice" {
				for i, s := range item[key].([]interface{}) {
					if s != nil && reflect.TypeOf(s).Kind().String() == "map" {
						item[key].([]interface{})[i] = normalizeMultiLookupsMap(s.(map[string]interface{}))
					}
				}
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

// For the use with defer to prevent static code
func shut(closer io.Closer) {
	if closer != nil {
		_ = closer.Close()
		// logs might be added here
	}
}
