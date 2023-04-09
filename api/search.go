package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/koltyakov/gosip"
)

// Search represents SharePoint Search API object struct
// Always use NewSearch constructor instead of &Search{}
type Search struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers *ODataMods
}

// SearchResp - search response type with helper processor methods
type SearchResp []byte

// SearchQuery - strongly typed struct for search method parameters
type SearchQuery struct {
	QueryText                             string                  `json:"Querytext"`                             // A string that contains the text for the search query
	QueryTemplate                         string                  `json:"QueryTemplate"`                         // A string that contains the text that replaces the query text, as part of a query transform
	EnableInterleaving                    bool                    `json:"EnableInterleaving"`                    // A Boolean value that specifies whether the result tables that are returned for the result block are mixed with the result tables that are returned for the original query
	EnableStemming                        bool                    `json:"EnableStemming"`                        // A Boolean value that specifies whether stemming is enabled
	TrimDuplicates                        bool                    `json:"TrimDuplicates"`                        // A Boolean value that specifies whether duplicate items are removed from the results
	EnableNicknames                       bool                    `json:"EnableNicknames"`                       // A Boolean value that specifies whether the exact terms in the search query are used to find matches, or if nicknames are used also
	EnableFQL                             bool                    `json:"EnableFQL"`                             // A Boolean value that specifies whether the query uses the FAST Query Language (FQL)
	EnablePhonetic                        bool                    `json:"EnablePhonetic"`                        // A Boolean value that specifies whether the phonetic forms of the query terms are used to find matches
	BypassResultTypes                     bool                    `json:"BypassResultTypes"`                     // A Boolean value that specifies whether to perform result type processing for the query
	ProcessBestBets                       bool                    `json:"ProcessBestBets"`                       // A Boolean value that specifies whether to return best bet results for the query. This parameter is used only when EnableQueryRules is set to true, otherwise it is ignored.
	EnableQueryRules                      bool                    `json:"EnableQueryRules"`                      // A Boolean value that specifies whether to enable query rules for the query
	EnableSorting                         bool                    `json:"EnableSorting"`                         // A Boolean value that specifies whether to sort search results
	GenerateBlockRankLog                  bool                    `json:"GenerateBlockRankLog"`                  // Specifies whether to return block rank log information in the BlockRankLog property of the interleaved result table. A block rank log contains the textual information on the block score and the documents that were de-duplicated.
	SourceID                              string                  `json:"SourceId"`                              // The result source ID to use for executing the search query
	RankingModelID                        string                  `json:"RankingModelId"`                        // The ID of the ranking model to use for the query
	StartRow                              int                     `json:"StartRow"`                              // The first row that is included in the search results that are returned. You use this parameter when you want to implement paging for search results.
	RowLimit                              int                     `json:"RowLimit"`                              // The maximum number of rows overall that are returned in the search results. Compared to RowsPerPage, RowLimit is the maximum number of rows returned overall.
	RowsPerPage                           int                     `json:"RowsPerPage"`                           // The maximum number of rows to return per page. Compared to RowLimit, RowsPerPage refers to the maximum number of rows to return per page, and is used primarily when you want to implement paging for search results.
	SelectProperties                      []string                `json:"SelectProperties"`                      // The managed properties to return in the search results
	Culture                               int                     `json:"Culture"`                               // The locale ID (LCID) for the query
	RefinementFilters                     []string                `json:"RefinementFilters"`                     // The set of refinement filters used when issuing a refinement query (FQL)
	Refiners                              string                  `json:"Refiners"`                              // The set of refiners to return in a search result
	HiddenConstraints                     string                  `json:"HiddenConstraints"`                     // The additional query terms to append to the query
	Timeout                               int                     `json:"Timeout"`                               // The amount of time in milliseconds before the query request times out
	HitHighlightedProperties              []string                `json:"HitHighlightedProperties"`              // The properties to highlight in the search result summary when the property value matches the search terms entered by the user
	ClientType                            string                  `json:"ClientType"`                            // The type of the client that issued the query
	PersonalizationData                   string                  `json:"PersonalizationData"`                   // The GUID for the user who submitted the search query
	ResultsURL                            string                  `json:"ResultsUrl"`                            // The URL for the search results page
	QueryTag                              string                  `json:"QueryTag"`                              // Custom tags that identify the query. You can specify multiple query tags
	ProcessPersonalFavorites              bool                    `json:"ProcessPersonalFavorites"`              // A Boolean value that specifies whether to return personal favorites with the search results
	QueryTemplatePropertiesURL            string                  `json:"QueryTemplatePropertiesUrl"`            // The location of the queryparametertemplate.xml file. This file is used to enable anonymous users to make Search REST queries
	HitHighlightedMultivaluePropertyLimit int                     `json:"HitHighlightedMultivaluePropertyLimit"` // The number of properties to show hit highlighting for in the search results
	EnableOrderingHitHighlightedProperty  bool                    `json:"EnableOrderingHitHighlightedProperty"`  // A Boolean value that specifies whether the hit highlighted properties can be ordered
	CollapseSpecification                 string                  `json:"CollapseSpecification"`                 // The managed properties that are used to determine how to collapse individual search results. Results are collapsed into one or a specified number of results if they match any of the individual collapse specifications. In a collapse specification, results are collapsed if their properties match all individual properties in the collapse specification.
	UIlanguage                            int                     `json:"UIlanguage"`                            // The locale identifier (LCID) of the user interface
	DesiredSnippetLength                  int                     `json:"DesiredSnippetLength"`                  // The preferred number of characters to display in the hit-highlighted summary generated for a search result
	MaxSnippetLength                      int                     `json:"MaxSnippetLength"`                      // The maximum number of characters to display in the hit-highlighted summary generated for a search result
	SummaryLength                         int                     `json:"SummaryLength"`                         // The number of characters to display in the result summary for a search result
	SortList                              []*SearchSort           `json:"SortList"`                              // The list of properties by which the search results are ordered
	Properties                            []*SearchProperty       `json:"Properties"`                            // Properties to be used to configure the search query
	ReorderingRules                       []*SearchReorderingRule `json:"ReorderingRules"`                       // Special rules for reordering search results. These rules can specify that documents matching certain conditions are ranked higher or lower in the results. This property applies only when search results are sorted based on rank.
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

// ResultTableCollection - search results table collection type
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
		client:    client,
		endpoint:  endpoint,
		config:    config,
		modifiers: NewODataMods(),
	}
}

func (search *Search) GetQuery(query *SearchQuery) (SearchResp, error) {
	sortList := query.SortList[0].Property + ":" + strconv.Itoa(query.SortList[0].Direction)
	endpoint := fmt.Sprintf("%s/query?querytext='%s'&sortlist='%s'&enabledsorting=true&rowlimit=%d",
		search.endpoint, query.QueryText, sortList, query.RowLimit)

	client := NewHTTPClient(search.client)

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json;odata=verbose;charset=utf-8",
	}

	return client.Get(endpoint, patchConfigHeaders(search.config, headers))
}

// PostQuery gets search results based on a `query`
func (search *Search) PostQuery(query *SearchQuery) (SearchResp, error) {
	endpoint := fmt.Sprintf("%s/PostQuery", search.endpoint)
	client := NewHTTPClient(search.client)

	request := map[string]interface{}{}
	queryBytes, _ := json.Marshal(query)
	_ = json.Unmarshal(queryBytes, &request)

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
	JSONReq := string(req)
	body := []byte(TrimMultiline(`{ "request": ` + JSONReq + `}`))

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json;odata=verbose;charset=utf-8",
	}

	return client.Post(endpoint, bytes.NewBuffer(body), patchConfigHeaders(search.config, headers))
}

// ToDo:
// _api/SP.UI.ApplicationPages.ClientPeoplePickerWebServiceInterface.ClientPeoplePickerSearchUser

/* Response helpers */

// Data : to get typed data
func (searchResp *SearchResp) Data() *SearchResults {
	data := NormalizeODataItem(*searchResp)
	res := &SearchResults{}
	_ = json.Unmarshal(data, &res)
	return res
}

// Normalized returns normalized body
func (searchResp *SearchResp) Normalized() []byte {
	return NormalizeODataItem(*searchResp)
}

// Results : to get typed data
func (searchResp *SearchResp) Results() []map[string]string {
	var results []map[string]string
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
