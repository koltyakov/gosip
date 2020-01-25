// Package api :: This is auto generated file, do not edit manually
package api

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (webs *Webs) Conf(config *RequestConfig) *Webs {
	webs.config = config
	return webs
}

// Select adds $select OData modifier
func (webs *Webs) Select(oDataSelect string) *Webs {
	webs.modifiers.AddSelect(oDataSelect)
	return webs
}

// Expand adds $expand OData modifier
func (webs *Webs) Expand(oDataExpand string) *Webs {
	webs.modifiers.AddExpand(oDataExpand)
	return webs
}

// Filter adds $filter OData modifier
func (webs *Webs) Filter(oDataFilter string) *Webs {
	webs.modifiers.AddFilter(oDataFilter)
	return webs
}

// Top adds $top OData modifier
func (webs *Webs) Top(oDataTop int) *Webs {
	webs.modifiers.AddTop(oDataTop)
	return webs
}

// OrderBy adds $orderby OData modifier
func (webs *Webs) OrderBy(oDataOrderBy string, ascending bool) *Webs {
	webs.modifiers.AddOrderBy(oDataOrderBy, ascending)
	return webs
}

/* Response helpers */

// Data response helper
func (websResp *WebsResp) Data() []WebResp {
	collection, _ := normalizeODataCollection(*websResp)
	webs := []WebResp{}
	for _, item := range collection {
		webs = append(webs, WebResp(item))
	}
	return webs
}

// Normalized returns normalized body
func (websResp *WebsResp) Normalized() []byte {
	normalized, _ := NormalizeODataCollection(*websResp)
	return normalized
}
