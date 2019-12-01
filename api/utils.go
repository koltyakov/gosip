package api

import (
	"encoding/json"
	"fmt"
	"net/url"
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
