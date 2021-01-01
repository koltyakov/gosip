package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/csom"
)

// Taxonomy session struct
type Taxonomy struct {
	client   *HTTPClient
	config   *RequestConfig
	endpoint string
}

// TermStore term store struct
type TermStore struct {
	client   *HTTPClient
	config   *RequestConfig
	endpoint string

	id   string
	name string

	selectProps []string
}

// TermGroups term groups struct
type TermGroups struct {
	client   *HTTPClient
	config   *RequestConfig
	endpoint string

	csomEntry csom.Builder
	termStore *TermStore
}

// TermGroup term group struct
type TermGroup struct {
	client   *HTTPClient
	config   *RequestConfig
	endpoint string

	id string

	csomEntry   csom.Builder
	selectProps []string
}

// TermSets term sets struct
type TermSets struct {
	client   *HTTPClient
	config   *RequestConfig
	endpoint string

	name string

	csomEntry csom.Builder
	termGroup *TermGroup
}

// TermSet term set struct
type TermSet struct {
	client   *HTTPClient
	config   *RequestConfig
	endpoint string

	id string

	csomEntry   csom.Builder
	selectProps []string
}

// Terms terms struct
type Terms struct {
	client   *HTTPClient
	config   *RequestConfig
	endpoint string

	csomEntry   csom.Builder
	termSet     *TermSet
	selectProps []string
}

// NewTaxonomy - taxonomy struct constructor function
func NewTaxonomy(client *gosip.SPClient, siteURL string, config *RequestConfig) *Taxonomy {
	return &Taxonomy{
		client:   NewHTTPClient(client),
		endpoint: siteURL,
		config:   config,
	}
}

// TermStore gets default site collection term store object
func (taxonomy *Taxonomy) TermStore() *TermStore {
	return &TermStore{
		client:   taxonomy.client,
		endpoint: taxonomy.endpoint,
		config:   taxonomy.config,

		id:   "",
		name: "",

		selectProps: []string{},
	}
}

// TermStoreByID gets term store object by ID
func (taxonomy *Taxonomy) TermStoreByID(termStoreGUID string) *TermStore {
	return &TermStore{
		client:   taxonomy.client,
		endpoint: taxonomy.endpoint,
		config:   taxonomy.config,

		id:   trimGUID(termStoreGUID),
		name: "",

		selectProps: []string{},
	}
}

// TermStoreByName gets term store object by Name
func (taxonomy *Taxonomy) TermStoreByName(termStoreName string) *TermStore {
	return &TermStore{
		client:   taxonomy.client,
		endpoint: taxonomy.endpoint,
		config:   taxonomy.config,

		id:   "",
		name: termStoreName,

		selectProps: []string{},
	}
}

// getCSOMBuilderEntry gets CSOM builder entry
func (termStore *TermStore) getCSOMBuilderEntry() csom.Builder {
	b := csom.NewBuilder()
	b.AddObject(csom.NewObject(`<StaticMethod TypeId="{981cbc68-9edc-4f8d-872f-71146fcbb84f}" Name="GetTaxonomySession" Id="{{.ID}}" />`), nil)
	if len(termStore.id) > 0 {
		// Term store by ID
		b.AddObject(csom.NewObjectProperty("TermStores"), nil)
		b.AddObject(csom.NewObjectMethod("GetById", []string{
			fmt.Sprintf(`<Parameter Type="String">%s</Parameter>`, termStore.id),
		}), nil)
	} else if len(termStore.name) > 0 {
		// Term store by Name
		b.AddObject(csom.NewObjectProperty("TermStores"), nil)
		b.AddObject(csom.NewObjectMethod("GetByName", []string{
			fmt.Sprintf(`<Parameter Type="String">%s</Parameter>`, termStore.name),
		}), nil)
	} else {
		// Default term store
		b.AddObject(csom.NewObjectMethod("GetDefaultSiteCollectionTermStore", []string{}), nil)
	}
	return b
}

// Select adds select props to term store query
func (termStore *TermStore) Select(props string) *TermStore {
	termStore.selectProps = appendProp(termStore.selectProps, props)
	return termStore
}

// Get gets term store metadata
func (termStore *TermStore) Get() (map[string]interface{}, error) {
	var props []string
	for _, prop := range termStore.selectProps {
		propertyXML := prop
		if strings.Index(prop, "<") == -1 {
			propertyXML = fmt.Sprintf(`<Property Name="%s" SelectAll="true" />`, prop)
		}
		props = append(props, propertyXML)
	}

	b := termStore.getCSOMBuilderEntry()
	b.AddAction(csom.NewQueryWithProps(props), nil)

	csomPkg, err := b.Compile()
	if err != nil {
		return nil, err
	}

	jsomResp, err := termStore.client.ProcessQuery(termStore.endpoint, bytes.NewBuffer([]byte(csomPkg)), termStore.config)
	if err != nil {
		return nil, err
	}

	var jsomRespArr []interface{}
	if err := json.Unmarshal(jsomResp, &jsomRespArr); err != nil {
		return nil, err
	}

	res, ok := jsomRespArr[len(jsomRespArr)-1].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("can't cast response, %v+", jsomRespArr[len(jsomRespArr)-1])
	}

	return res, nil
}

// Groups gets term groups object
func (termStore *TermStore) Groups() *TermGroups {
	return &TermGroups{
		client:    termStore.client,
		endpoint:  termStore.endpoint,
		config:    termStore.config,
		csomEntry: termStore.getCSOMBuilderEntry(),
		termStore: termStore,
	}
}

// Get gets term groups metadata
func (termGroups *TermGroups) Get() ([]map[string]interface{}, error) {
	termStore, err := termGroups.termStore.Select("Groups").Get()
	if err != nil {
		return nil, err
	}

	groups, ok := termStore["Groups"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("can't get groups from term store")
	}

	items, ok := groups["_Child_Items_"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("can't get child items from groups")
	}

	var resItems []map[string]interface{}
	for _, item := range items {
		resItem, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("can't get item from groups")
		}
		resItems = append(resItems, resItem)
	}

	return resItems, nil
}

// GetByID gets term group object by its GUID
func (termGroups *TermGroups) GetByID(groupGUID string) *TermGroup {
	return &TermGroup{
		client:   termGroups.client,
		endpoint: termGroups.endpoint,
		config:   termGroups.config,
		id:       trimGUID(groupGUID),

		csomEntry:   termGroups.csomEntry.Clone(),
		selectProps: []string{},
	}
}

// getCSOMBuilderEntry gets CSOM builder entry
func (termGroup *TermGroup) getCSOMBuilderEntry() csom.Builder {
	b := termGroup.csomEntry.Clone()
	b.AddObject(csom.NewObjectMethod("GetGroup", []string{
		fmt.Sprintf(`<Parameter Type="String">%s</Parameter>`, termGroup.id),
	}), nil)
	return b
}

// Select adds select props to term store query
func (termGroup *TermGroup) Select(props string) *TermGroup {
	termGroup.selectProps = appendProp(termGroup.selectProps, props)
	return termGroup
}

// Get gets term group metadata
func (termGroup *TermGroup) Get() (map[string]interface{}, error) {
	var props []string
	for _, prop := range termGroup.selectProps {
		propertyXML := prop
		if strings.Index(prop, "<") == -1 {
			propertyXML = fmt.Sprintf(`<Property Name="%s" SelectAll="true" />`, prop)
		}
		props = append(props, propertyXML)
	}

	b := termGroup.getCSOMBuilderEntry()
	b.AddAction(csom.NewQueryWithProps(props), nil)

	csomPkg, err := b.Compile()
	if err != nil {
		return nil, err
	}

	jsomResp, err := termGroup.client.ProcessQuery(termGroup.endpoint, bytes.NewBuffer([]byte(csomPkg)), termGroup.config)
	if err != nil {
		return nil, err
	}

	var jsomRespArr []interface{}
	if err := json.Unmarshal(jsomResp, &jsomRespArr); err != nil {
		return nil, err
	}

	res, ok := jsomRespArr[len(jsomRespArr)-1].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("can't cast response, %v+", jsomRespArr[len(jsomRespArr)-1])
	}

	return res, nil
}

// TermSets gets term sets object for current term group
func (termGroup *TermGroup) TermSets() *TermSets {
	return &TermSets{
		client:    termGroup.client,
		endpoint:  termGroup.endpoint,
		config:    termGroup.config,
		csomEntry: termGroup.getCSOMBuilderEntry(),
		termGroup: termGroup,
	}
}

// Get gets term sets metadata
func (termSets *TermSets) Get() ([]map[string]interface{}, error) {
	termStore, err := termSets.termGroup.Select("TermSets").Get()
	if err != nil {
		return nil, err
	}

	groups, ok := termStore["TermSets"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("can't get term sets from term group")
	}

	items, ok := groups["_Child_Items_"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("can't get child items from term sets")
	}

	var resItems []map[string]interface{}
	for _, item := range items {
		resItem, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("can't get item from term set")
		}
		resItems = append(resItems, resItem)
	}

	return resItems, nil
}

// GetTermSet gets term set object by its GUID
func (termStore *TermStore) GetTermSet(termSetGUID string) *TermSet {
	return &TermSet{
		client:   termStore.client,
		endpoint: termStore.endpoint,
		config:   termStore.config,

		id: trimGUID(termSetGUID),

		csomEntry:   termStore.getCSOMBuilderEntry().Clone(),
		selectProps: []string{},
	}
}

// getCSOMBuilderEntry gets CSOM builder entry
func (termSet *TermSet) getCSOMBuilderEntry() csom.Builder {
	b := termSet.csomEntry.Clone()
	b.AddObject(csom.NewObjectMethod("GetTermSet", []string{
		fmt.Sprintf(`<Parameter Type="String">%s</Parameter>`, termSet.id),
	}), nil)
	return b
}

// Select adds select props to term set query
func (termSet *TermSet) Select(props string) *TermSet {
	termSet.selectProps = appendProp(termSet.selectProps, props)
	return termSet
}

// Get gets term set metadata
func (termSet *TermSet) Get() (map[string]interface{}, error) {
	var props []string
	for _, prop := range termSet.selectProps {
		propertyXML := prop
		if strings.Index(prop, "<") == -1 {
			propertyXML = fmt.Sprintf(`<Property Name="%s" SelectAll="true" />`, prop)
		}
		props = append(props, propertyXML)
	}

	b := termSet.getCSOMBuilderEntry().Clone()
	b.AddAction(csom.NewQueryWithProps(props), nil)

	csomPkg, err := b.Compile()
	if err != nil {
		return nil, err
	}

	jsomResp, err := termSet.client.ProcessQuery(termSet.endpoint, bytes.NewBuffer([]byte(csomPkg)), termSet.config)
	if err != nil {
		return nil, err
	}

	var jsomRespArr []interface{}
	if err := json.Unmarshal(jsomResp, &jsomRespArr); err != nil {
		return nil, err
	}

	res, ok := jsomRespArr[len(jsomRespArr)-1].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("can't cast response, %v+", jsomRespArr[len(jsomRespArr)-1])
	}

	return res, nil
}

// Terms gets terms object instance
func (termSet *TermSet) Terms() *Terms {
	return &Terms{
		client:   termSet.client,
		endpoint: termSet.endpoint,
		config:   termSet.config,

		csomEntry:   termSet.getCSOMBuilderEntry().Clone(),
		termSet:     termSet,
		selectProps: []string{},
	}
}

// Select adds select props to terms collection query
func (terms *Terms) Select(props string) *Terms {
	terms.selectProps = appendProp(terms.selectProps, props)
	return terms
}

// GetAll gets all terms collection metadata
func (terms *Terms) GetAll() ([]map[string]interface{}, error) {
	var props []string
	for _, prop := range terms.selectProps {
		propertyXML := prop
		if strings.Index(prop, "<") == -1 {
			propertyXML = fmt.Sprintf(`<Property Name="%s" SelectAll="true" />`, prop)
		}
		props = append(props, propertyXML)
	}

	b := terms.termSet.getCSOMBuilderEntry().Clone()
	b.AddObject(csom.NewObjectMethod("GetAllTerms", []string{}), nil)
	b.AddAction(csom.NewQueryWithChildProps(props), nil)

	csomPkg, err := b.Compile()
	if err != nil {
		return nil, err
	}

	jsomResp, err := terms.client.ProcessQuery(terms.endpoint, bytes.NewBuffer([]byte(csomPkg)), terms.config)
	if err != nil {
		return nil, err
	}

	var jsomRespArr []interface{}
	if err := json.Unmarshal(jsomResp, &jsomRespArr); err != nil {
		return nil, err
	}

	termsResp, ok := jsomRespArr[len(jsomRespArr)-1].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("can't cast response, %v+", jsomRespArr[len(jsomRespArr)-1])
	}

	items, ok := termsResp["_Child_Items_"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("can't get child items from groups")
	}

	var resItems []map[string]interface{}
	for _, item := range items {
		resItem, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("can't get item from groups")
		}
		resItems = append(resItems, resItem)
	}

	return resItems, nil
}

// Utility methods

func appendProp(props []string, prop string) []string {
	for _, p := range strings.SplitN(prop, ",", -1) {
		p = strings.Trim(p, " ")
		found := false
		for _, pp := range props {
			if pp == p {
				found = true
			}
		}
		if !found {
			props = append(props, p)
		}
	}
	return props
}

func trimGUID(guid string) string {
	guid = strings.ToLower(guid)
	guid = strings.Replace(guid, "/guid(", "", 1)
	guid = strings.Replace(guid, ")/", "", 1)
	return guid
}
