package api

import (
	"fmt"
	"net/url"
)

// ODataMods - REST OData Modifiers struct
type ODataMods struct {
	mods map[string]string
}

// NewODataMods - ODataMods constructor function
func NewODataMods() *ODataMods {
	return &ODataMods{
		mods: map[string]string{},
	}
}

// Get retrieves OData modifiers map
func (oData *ODataMods) Get() map[string]string {
	if oData == nil {
		oData = &ODataMods{mods: map[string]string{}}
	}
	if oData.mods == nil {
		oData.mods = map[string]string{}
	}
	return oData.mods
}

// AddSelect adds $select OData modifier
func (oData *ODataMods) AddSelect(values string) *ODataMods {
	if oData == nil {
		oData = &ODataMods{mods: map[string]string{}}
	}
	if oData.mods == nil {
		oData.mods = map[string]string{}
	}
	oData.mods["$select"] = values
	return oData
}

// AddExpand adds $expand OData modifier
func (oData *ODataMods) AddExpand(values string) *ODataMods {
	if oData == nil {
		oData = &ODataMods{mods: map[string]string{}}
	}
	if oData.mods == nil {
		oData.mods = map[string]string{}
	}
	oData.mods["$expand"] = values
	return oData
}

// AddFilter adds $filter OData modifier
func (oData *ODataMods) AddFilter(values string) *ODataMods {
	if oData == nil {
		oData = &ODataMods{mods: map[string]string{}}
	}
	if oData.mods == nil {
		oData.mods = map[string]string{}
	}
	oData.mods["$filter"] = values
	return oData
}

// AddSkip adds $skiptoken OData modifier
func (oData *ODataMods) AddSkip(value string) *ODataMods {
	if oData == nil {
		oData = &ODataMods{mods: map[string]string{}}
	}
	if oData.mods == nil {
		oData.mods = map[string]string{}
	}
	oData.mods["$skiptoken"] = value
	return oData
}

// AddTop adds $top OData modifier
func (oData *ODataMods) AddTop(value int) *ODataMods {
	if oData == nil {
		oData = &ODataMods{mods: map[string]string{}}
	}
	if oData.mods == nil {
		oData.mods = map[string]string{}
	}
	oData.mods["$top"] = fmt.Sprintf("%d", value)
	return oData
}

// AddOrderBy adds $orderby OData modifier
func (oData *ODataMods) AddOrderBy(orderBy string, ascending bool) *ODataMods {
	if oData == nil {
		oData = &ODataMods{mods: map[string]string{}}
	}
	if oData.mods == nil {
		oData.mods = map[string]string{}
	}
	direction := "asc"
	if !ascending {
		direction = "desc"
	}
	if oData.mods["$orderby"] != "" {
		oData.mods["$orderby"] += ","
	}
	oData.mods["$orderby"] += fmt.Sprintf("%s %s", orderBy, direction)
	return oData
}

// Endpoint with OData modifiers toURL helper method
func toURL(endpoint string, modifiers *ODataMods) string {
	apiURL, _ := url.Parse(endpoint)
	query := apiURL.Query() // url.Values{}
	for k, v := range modifiers.Get() {
		query.Set(k, trimMultiline(v))
	}
	apiURL.RawQuery = query.Encode()
	return apiURL.String()
}
