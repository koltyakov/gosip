package api

// StringValue ...
type StringValue struct {
	StringValue string `json:"StringValue"`
}

// DecodedURL ...
type DecodedURL struct {
	DecodedURL string `json:"DecodedUrl"`
}

// TypedKeyValue - typed key value prop
type TypedKeyValue struct {
	Key       string `json:"Key"`
	Value     string `json:"Value"`
	ValueType string `json:"ValueType"`
}
