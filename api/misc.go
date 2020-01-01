package api

// StringValue single value prop type
type StringValue struct {
	StringValue string `json:"StringValue"`
}

// DecodedURL decode URL prop type
type DecodedURL struct {
	DecodedURL string `json:"DecodedUrl"`
}

// TypedKeyValue typed key value prop type
type TypedKeyValue struct {
	Key       string `json:"Key"`
	Value     string `json:"Value"`
	ValueType string `json:"ValueType"`
}
