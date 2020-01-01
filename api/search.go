package api

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/koltyakov/gosip"
)

// Search represents SharePoint Search API object struct
type Search struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers map[string]string
}

// SearchResp - search response type with helper processor methods
type SearchResp []byte

// SearchQuery - strongly typed struct for search method parameters
type SearchQuery struct {
	QueryText                             string                  `json:"Querytext"`
	QueryTemplate                         string                  `json:"QueryTemplate"`
	EnableInterleaving                    bool                    `json:"EnableInterleaving"`
	EnableStemming                        bool                    `json:"EnableStemming"`
	TrimDuplicates                        bool                    `json:"TrimDuplicates"`
	EnableNicknames                       bool                    `json:"EnableNicknames"`
	EnableFQL                             bool                    `json:"EnableFQL"`
	EnablePhonetic                        bool                    `json:"EnablePhonetic"`
	BypassResultTypes                     bool                    `json:"BypassResultTypes"`
	ProcessBestBets                       bool                    `json:"ProcessBestBets"`
	EnableQueryRules                      bool                    `json:"EnableQueryRules"`
	EnableSorting                         bool                    `json:"EnableSorting"`
	GenerateBlockRankLog                  bool                    `json:"GenerateBlockRankLog"`
	SourceID                              string                  `json:"SourceId"`
	RankingModelID                        string                  `json:"RankingModelId"`
	StartRow                              int                     `json:"StartRow"`
	RowLimit                              int                     `json:"RowLimit"`
	RowsPerPage                           int                     `json:"RowsPerPage"`
	SelectProperties                      []string                `json:"SelectProperties"`
	Culture                               int                     `json:"Culture"`
	RefinementFilters                     []string                `json:"RefinementFilters"`
	Refiners                              string                  `json:"Refiners"`
	HiddenConstraints                     string                  `json:"HiddenConstraints"`
	Timeout                               int                     `json:"Timeout"`
	HitHighlightedProperties              []string                `json:"HitHighlightedProperties"`
	ClientType                            string                  `json:"ClientType"`
	PersonalizationData                   string                  `json:"PersonalizationData"`
	ResultsURL                            string                  `json:"ResultsUrl"`
	QueryTag                              string                  `json:"QueryTag"`
	ProcessPersonalFavorites              bool                    `json:"ProcessPersonalFavorites"`
	QueryTemplatePropertiesURL            string                  `json:"QueryTemplatePropertiesUrl"`
	HitHighlightedMultivaluePropertyLimit int                     `json:"HitHighlightedMultivaluePropertyLimit"`
	EnableOrderingHitHighlightedProperty  bool                    `json:"EnableOrderingHitHighlightedProperty"`
	CollapseSpecification                 string                  `json:"CollapseSpecification"`
	UIlanguage                            int                     `json:"UIlanguage"`
	DesiredSnippetLength                  int                     `json:"DesiredSnippetLength"`
	MaxSnippetLength                      int                     `json:"MaxSnippetLength"`
	SummaryLength                         int                     `json:"SummaryLength"`
	SortList                              []*SearchSort           `json:"SortList"`
	Properties                            []*SearchProperty       `json:"Properties"`
	ReorderingRules                       []*SearchReorderingRule `json:"ReorderingRules"`
}

// SearchSort - search sort property type
type SearchSort struct {
	Property  string `json:"Property"`
	Direction int    `json:"Direction"` // Ascending = 0, Descending = 1, FQLFormula = 2
}

// SearchProperty - search property type
type SearchProperty struct {
	Name  string              `json:"Name"`
	Value SearchPropertyValue `json:"Value"`
}

// SearchPropertyValue - search property value type
type SearchPropertyValue struct {
	StrVal                      string   `json:"StrVal"`
	BoolVal                     bool     `json:"BoolVal"`
	IntVal                      int      `json:"IntVal"`
	StrArray                    []string `json:"StrArray"`
	QueryPropertyValueTypeIndex int      `json:"QueryPropertyValueTypeIndex"` // None = 0, StringType = 1, Int32Type = 2, BooleanType = 3, StringArrayType = 4, UnSupportedType = 5
}

// SearchReorderingRule - search reordering rule type
type SearchReorderingRule struct {
	MatchValue string `json:"MatchValue"`
	Boost      int    `json:"Boost"`
	MatchType  int    `json:"MatchType"` // ResultContainsKeyword = 0, TitleContainsKeyword = 1, TitleMatchesKeyword = 2, UrlStartsWith = 3, UrlExactlyMatches = 4, ContentTypeIs = 5, FileExtensionMatches = 6, ResultHasTag = 7, ManualCondition = 8
}

// SearchResults - search results response type
type SearchResults struct {
	ElapsedTime           int                      `json:"ElapsedTime"`
	PrimaryQueryResult    *ResultTableCollection   `json:"PrimaryQueryResult"`
	Properties            []*TypedKeyValue         `json:"Properties"`
	SecondaryQueryResults []*ResultTableCollection `json:"SecondaryQueryResults"`
	SpellingSuggestion    string                   `json:"SpellingSuggestion"`
	TriggeredRules        []interface{}            `json:"TriggeredRules"`
}

// ResultTableCollection - search results table collecton type
type ResultTableCollection struct {
	QueryErrors        map[string]interface{} `json:"QueryErrors"`
	QueryID            string                 `json:"QueryId"`
	QueryRuleID        string                 `json:"QueryRuleId"`
	CustomResults      *ResultTable           `json:"CustomResults"`
	RefinementResults  *ResultTable           `json:"RefinementResults"`
	RelevantResults    *ResultTable           `json:"RelevantResults"`
	SpecialTermResults *ResultTable           `json:"SpecialTermResults"`
}

// ResultTable - search result table type
type ResultTable struct {
	GroupTemplateID              string           `json:"GroupTemplateId"`
	ItemTemplateID               string           `json:"ItemTemplateId"`
	ResultTitle                  string           `json:"ResultTitle"`
	ResultTitleURL               string           `json:"ResultTitleUrl"`
	RowCount                     int              `json:"RowCount"`
	TableType                    string           `json:"TableType"`
	TotalRows                    int              `json:"TotalRows"`
	TotalRowsIncludingDuplicates int              `json:"TotalRowsIncludingDuplicates"`
	Properties                   []*TypedKeyValue `json:"Properties"`
	Table                        *struct {
		Rows []*struct {
			Cells []*TypedKeyValue `json:"Cells"`
		} `json:"Rows"`
	} `json:"Table"`
	Refiners []*struct {
		Name    string `json:"Name"`
		Entries []*struct {
			RefinementCount string `json:"RefinementCount"`
			RefinementName  string `json:"RefinementName"`
			RefinementToken string `json:"RefinementToken"`
			RefinementValue string `json:"RefinementValue"`
		} `json:"Entries"`
	} `json:"Refiners"`
}

// NewSearch - Search struct constructor function
func NewSearch(client *gosip.SPClient, endpoint string, config *RequestConfig) *Search {
	return &Search{
		client:   client,
		endpoint: endpoint,
		config:   config,
	}
}

// PostQuery ...
func (search *Search) PostQuery(query *SearchQuery) (SearchResp, error) {
	endpoint := fmt.Sprintf("%s/PostQuery", search.endpoint)
	sp := NewHTTPClient(search.client)

	request := map[string]interface{}{}
	queryBytes, _ := json.Marshal(query)
	json.Unmarshal(queryBytes, &request)

	request["__metadata"] = map[string]interface{}{"type": "Microsoft.Office.Server.Search.REST.SearchRequest"}

	for key, val := range request {
		if val == nil {
			delete(request, key)
		} else {
			switch kind := reflect.TypeOf(val).Kind().String(); kind {
			case "slice":
				request[key] = map[string]interface{}{
					"results": val,
				}
			case "int":
				if val.(int) == 0 {
					delete(request, key)
				}
			case "float64":
				if val.(float64) == 0 {
					delete(request, key)
				}
			case "string":
				if val.(string) == "" {
					delete(request, key)
				}
			default:
			}
		}
	}

	req, _ := json.Marshal(request)
	JSONReq := fmt.Sprintf("%s", req)
	body := []byte(trimMultiline(`{ "request": ` + JSONReq + `}`))

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json;odata=verbose;charset=utf-8",
	}

	return sp.Post(endpoint, body, headers)
}

// ToDo:
// _api/SP.UI.ApplicationPages.ClientPeoplePickerWebServiceInterface.ClientPeoplePickerSearchUser

/* Response helpers */

// Data : to get typed data
func (searchResp *SearchResp) Data() *SearchResults {
	data := parseODataItem(*searchResp)
	res := &SearchResults{}
	json.Unmarshal(data, &res)
	return res
}

// Results : to get typed data
func (searchResp *SearchResp) Results() []map[string]string {
	results := []map[string]string{}
	rows := searchResp.Data().PrimaryQueryResult.RelevantResults.Table.Rows
	for _, row := range rows {
		rowMap := map[string]string{}
		for _, cell := range row.Cells {
			rowMap[cell.Key] = cell.Value
		}
		results = append(results, rowMap)
	}
	return results
}
