package api

import (
	"fmt"
	"strings"

	"github.com/koltyakov/gosip/csom"
)

// TermStores term stores struct
type TermStores struct {
	client   *HTTPClient
	config   *RequestConfig
	endpoint string
}

// Default gets default site collection term store object
func (termStores *TermStores) Default() *TermStore {
	return &TermStore{
		client:   termStores.client,
		endpoint: termStores.endpoint,
		config:   termStores.config,

		id:   "",
		name: "",

		selectProps: []string{},
	}
}

// GetByID gets term store object by ID
func (termStores *TermStores) GetByID(storeGUID string) *TermStore {
	return &TermStore{
		client:   termStores.client,
		endpoint: termStores.endpoint,
		config:   termStores.config,

		id:   trimTaxonomyGUID(storeGUID),
		name: "",

		selectProps: []string{},
	}
}

// GetByName gets term store object by Name
func (termStores *TermStores) GetByName(storeName string) *TermStore {
	return &TermStore{
		client:   termStores.client,
		endpoint: termStores.endpoint,
		config:   termStores.config,

		id:   "",
		name: storeName,

		selectProps: []string{},
	}
}

/* Term store */

// TermStore term store struct
type TermStore struct {
	client   *HTTPClient
	config   *RequestConfig
	endpoint string

	id   string
	name string

	selectProps []string
}

// csomBuilderEntry gets CSOM builder entry
func (termStore *TermStore) csomBuilderEntry() csom.Builder {
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
	termStore.selectProps = appendTaxonomyProp(termStore.selectProps, props)
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

	b := termStore.csomBuilderEntry()
	b.AddAction(csom.NewQueryWithProps(props), nil)

	return getCSOMResponse(termStore.client, termStore.endpoint, termStore.config, b)
}

// // CommitAll commits all changes
// func (termStore *TermStore) CommitAll() error {
// 	b := termStore.csomBuilderEntry().Clone()
// 	b.AddAction(csom.NewActionMethod("CommitAll", []string{}), nil)
// 	_, err := getCSOMResponse(termStore.client, termStore.endpoint, termStore.config, b)
// 	return err
// }

// UpdateCache updates store cache
func (termStore *TermStore) UpdateCache() error {
	b := termStore.csomBuilderEntry().Clone()
	b.AddAction(csom.NewActionMethod("UpdateCache", []string{}), nil)
	_, err := getCSOMResponse(termStore.client, termStore.endpoint, termStore.config, b)
	return err
}

// Groups gets term groups object
func (termStore *TermStore) Groups() *TermGroups {
	return &TermGroups{
		client:    termStore.client,
		endpoint:  termStore.endpoint,
		config:    termStore.config,
		csomEntry: termStore.csomBuilderEntry(),
	}
}

// Sets gets term sets object
func (termStore *TermStore) Sets() *TermSets {
	return &TermSets{
		client:    termStore.client,
		endpoint:  termStore.endpoint,
		config:    termStore.config,
		csomEntry: termStore.csomBuilderEntry(),
	}
}

// Terms gets terms object
func (termStore *TermStore) Terms() *Terms {
	return &Terms{
		client:    termStore.client,
		endpoint:  termStore.endpoint,
		config:    termStore.config,
		csomEntry: termStore.csomBuilderEntry(),
	}
}
