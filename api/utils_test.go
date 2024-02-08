package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestUtils(t *testing.T) {

	t.Run("getConfHeaders", func(t *testing.T) {
		h := getConfHeaders(headers.verbose)
		if !reflect.DeepEqual(h, headers.verbose.Headers) {
			t.Errorf(
				"incorrect headers, expected \"%s\", got \"%s\"",
				headers.verbose.Headers,
				h,
			)
		}
	})

	t.Run("trimMultiline", func(t *testing.T) {
		initial := `
			<div>
				Hello SharePoint
			</div>
		`
		expected := "<div>Hello SharePoint</div>"
		resStr := TrimMultiline(initial)
		if resStr != expected {
			t.Errorf(
				"incorrect string, expected \"%s\", got \"%s\"",
				expected,
				resStr,
			)
		}
	})

	t.Run("getRelativeURL", func(t *testing.T) {
		absURL := "https://contoso.sharepoint.com/sites/site"
		relURL := "/sites/site"
		resRelURL := getRelativeURL(absURL)
		if resRelURL != relURL {
			t.Errorf(
				"incorrect relative URL, expected \"%s\", got \"%s\"",
				relURL,
				resRelURL,
			)
		}
	})

	t.Run("getPriorEndpoint", func(t *testing.T) {
		endpoint := "https://contoso.sharepoint.com/sites/site/_api/web/Lists('UpperCaseList')"
		part := "/Items"
		endpoint1 := endpoint + "/Items/Item(1)"
		endpoint2 := endpoint + "/items/item(1)"
		endpoint3 := endpoint + "/IDK/Item(1)"
		if getPriorEndpoint(endpoint1, part) != endpoint {
			t.Error("incorrect endpoint reduction")
		}
		if getPriorEndpoint(endpoint2, part) != endpoint {
			t.Error("incorrect endpoint reduction")
		}
		if getPriorEndpoint(endpoint3, part) == endpoint {
			t.Error("incorrect endpoint reduction")
		}
	})

	t.Run("getIncludeEndpoint", func(t *testing.T) {
		rootEndpoint := "https://contoso.sharepoint.com/sites/site/_api/web/Lists('UpperCaseList')"
		part := "/Items"
		endpoint := rootEndpoint + part + "/Item(1)"
		if getIncludeEndpoint(endpoint, part) != rootEndpoint+part {
			t.Error("incorrect endpoint reduction")
		}
	})

	t.Run("getIncludeEndpoints", func(t *testing.T) {
		rootEndpoint := "https://contoso.sharepoint.com/sites/site/_api/web/Lists('UpperCaseList')"
		part := "/Items"
		endpoint := rootEndpoint + part + "/Item(1)"
		parts := []string{"/Items", "/SomethingElse"}
		if getIncludeEndpoints(endpoint, parts) != rootEndpoint+part {
			t.Error("incorrect endpoint reduction")
		}
	})

	t.Run("containsMetadataType", func(t *testing.T) {
		m1 := []byte(`{"prop":"val"}`)
		if containsMetadataType(m1) {
			t.Error("metadata was detected but actually should not")
		}
		m2 := []byte(`{"__metadata":{"type":"SP.Any"},"prop":"val"}`)
		if !containsMetadataType(m2) {
			t.Error("metadata was not detected but actually should")
		}
	})

	t.Run("patchMetadataType", func(t *testing.T) {
		// should add __metadata if missed
		m1 := []byte(`{"prop":"val"}`)
		p1 := patchMetadataType(m1, "SP.List")
		if !containsMetadataType(p1) {
			t.Error("payload was not enriched with metadata")
		}
		if !strings.Contains(string(p1), `"prop":"val"`) {
			t.Error("payload metadata lost prop(s)")
		}

		// should not add __metadata if present
		m2 := []byte(`{"__metadata":{"type":"SP.Any"},"prop":"val"}`)
		p2 := patchMetadataType(m2, "SP.List")
		if !containsMetadataType(p2) {
			t.Error("payload lost metadata prop")
		}
		if !strings.Contains(string(p2), `"__metadata":{"type":"SP.Any"}`) {
			t.Error("payload metadata prop was mutated")
		}
	})

	t.Run("patchMetadataTypeCB", func(t *testing.T) {
		m1 := []byte(`{"prop":"val"}`)
		p1 := patchMetadataTypeCB(m1, func() string {
			return "SP.List"
		})
		if !containsMetadataType(p1) {
			t.Error("payload was not enriched with metadata")
		}
	})

	t.Run("parseODataItem", func(t *testing.T) {
		minimal := []byte(`{"prop":"val"}`)
		verbose := []byte(fmt.Sprintf(`{"d":%s}`, minimal))
		if !bytes.Contains(NormalizeODataItem(verbose), minimal) {
			t.Error("wrong OData transformation")
		}
		if !bytes.Contains(NormalizeODataItem(minimal), minimal) {
			t.Error("wrong OData transformation")
		}
	})

	t.Run("parseODataCollection", func(t *testing.T) {
		resulted := []byte(`[{"prop":"val1"},{"prop":"val2"}]`)
		verbose := []byte(fmt.Sprintf(`{"d":{"results":%s}}`, resulted))
		var fromVerbose []byte
		coll, _ := normalizeODataCollection(verbose)
		for _, b := range coll {
			fromVerbose = append(fromVerbose, b...)
		}
		var fromMinimal []byte
		coll, _ = normalizeODataCollection(resulted)
		for _, b := range coll {
			fromMinimal = append(fromMinimal, b...)
		}
		if !bytes.Equal(fromVerbose, fromMinimal) {
			t.Error("wrong OData transformation")
		}
	})

	t.Run("parseODataCollection2", func(t *testing.T) {
		resulted := []byte(`[{"prop":"val1"},{"prop":"val2"}]`)
		minimal := []byte(fmt.Sprintf(`{"value":%s}`, resulted))
		verbose := []byte(fmt.Sprintf(`{"d":{"results":%s}}`, resulted))
		var fromVerbose []byte
		coll, _ := normalizeODataCollection(verbose)
		for _, b := range coll {
			fromVerbose = append(fromVerbose, b...)
		}
		var fromMinimal []byte
		coll, _ = normalizeODataCollection(minimal)
		for _, b := range coll {
			fromMinimal = append(fromMinimal, b...)
		}
		if !bytes.Equal(fromVerbose, fromMinimal) {
			t.Error("wrong OData transformation")
		}
	})

	t.Run("parseODataCollection/Empty", func(t *testing.T) {
		minimal := []byte(`[]`)
		verbose := []byte(fmt.Sprintf(`{"d":{"results":%s}}`, minimal))
		var fromVerbose []byte
		coll, _ := normalizeODataCollection(verbose)
		for _, b := range coll {
			fromVerbose = append(fromVerbose, b...)
		}
		var fromMinimal []byte
		coll, _ = normalizeODataCollection(minimal)
		for _, b := range coll {
			fromMinimal = append(fromMinimal, b...)
		}
		if !bytes.Equal(fromVerbose, fromMinimal) {
			t.Error("wrong OData transformation")
		}
	})

	t.Run("normalizeMultiLookups", func(t *testing.T) {
		minimal := []byte(`{"multi":[1,2,3],"single":"val1"}`)
		verbose := []byte(`{"multi":{"results":[1,2,3]},"single":"val1"}`)
		if !bytes.Equal(normalizeMultiLookups(verbose), minimal) {
			t.Error("wrong OData transformation")
		}
	})

	t.Run("normalizeMultiLookupsMap", func(t *testing.T) {
		raw := []byte(`{"test1":{"results":[1,2,3]},"test2":{"other":[1,2,3]}}`)
		expected := []byte(`{"test1":[1,2,3],"test2":{"other":[1,2,3]}}`)

		rawMap := map[string]interface{}{}
		_ = json.Unmarshal(raw, &rawMap)
		resMap := normalizeMultiLookupsMap(rawMap)

		res, err := json.Marshal(resMap)
		if err != nil {
			t.Error(err)
		}

		if !bytes.Equal(res, expected) {
			t.Error("wrong transformation")
		}
	})

	t.Run("normalizeMultiLookupsMap/WithNulls", func(t *testing.T) {
		raw := []byte(`{
			"ID": 144,
			"SPFTaskCompeted": true,
			"SPFTaskDescription": null,
			"SPFTaskResult": "{\"completeAction\":\"RESEARCH_DONE\",\"routes\":[]}",
			"SPFTaskType": "STUDY",
			"__metadata": {
				"etag": "\"4\"",
				"id": "Web/Lists(guid'c0c8ab0e-2a5f-40cf-83e0-b7c0ec98eed6')/Items(144)",
				"type": "SP.Data.TasksListItem",
				"uri": "http://sp.contoso.com/sites/site/_api/Web/Lists(guid'c0c8ab0e-2a5f-40cf-83e0-b7c0ec98eed6')/Items(144)"
			}
		}`)

		rawMap := map[string]interface{}{}
		_ = json.Unmarshal(raw, &rawMap)
		resMap := normalizeMultiLookupsMap(rawMap)

		_, err := json.Marshal(resMap)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("fixDatesInResponse", func(t *testing.T) {
		data := []byte(`{
			"valid": "2019-12-03T12:19:45Z",
			"invalid": "2019-12-03T12:19:45"
		}`)
		d := fixDatesInResponse(data, []string{"valid", "invalid"})
		r := &struct {
			Valid   time.Time `json:"valid"`
			Invalid time.Time `json:"invalid"`
		}{}
		if err := json.Unmarshal(d, &r); err != nil {
			t.Error(err)
		}
		if r.Valid != r.Invalid {
			t.Error("dates have not been fixed in payload")
		}
	})

	t.Run("getODataCollectionNextPageURL", func(t *testing.T) {
		pageURL := "page-url"
		verbose := []byte(`{ "d": { "__next": "` + pageURL + `" }}`)
		minimal := []byte(`{ "odata.nextLink": "` + pageURL + `" }`)

		vNextPage := getODataCollectionNextPageURL(verbose)
		if vNextPage != pageURL {
			t.Error("can't parse next page")
		}

		mNextPage := getODataCollectionNextPageURL(minimal)
		if mNextPage != pageURL {
			t.Error("can't parse next page")
		}
	})

	t.Run("checkGetRelativeURL", func(t *testing.T) {
		ctxURL := "https://contoso.sharepoint.com/sites/site/my-site"
		relativeURI := "Shared Documents/Folder"
		resultURL := "/sites/site/my-site/Shared Documents/Folder"
		res := checkGetRelativeURL(relativeURI, ctxURL)
		if resultURL != res {
			t.Errorf(`wrong URL transformation, expected "%s", received "%s"`, resultURL, res)
		}
	})

	t.Run("checkGetRelativeURLEmpty", func(t *testing.T) {
		ctxURL := "https://contoso.sharepoint.com/sites/site/my-site"
		relativeURI := ""
		resultURL := "/sites/site/my-site"
		res := checkGetRelativeURL(relativeURI, ctxURL)
		if resultURL != res {
			t.Errorf(`wrong URL transformation, expected "%s", received "%s"`, resultURL, res)
		}
	})

	t.Run("extractEntityURI", func(t *testing.T) {
		ep1 := []byte(`{
			"__metadata": {
				"id": "entity_url"
			}
		}`)
		if ExtractEntityURI(ep1) != "entity_url" {
			t.Error("can't extract entity URL")
		}

		ep2 := []byte(`{
			"odata.id": "entity_url"
		}`)
		if ExtractEntityURI(ep2) != "entity_url" {
			t.Error("can't extract entity URL")
		}
	})

	t.Run("patchConfigHeaders", func(t *testing.T) {
		headers := map[string]string{
			"Accept": "application/json",
		}
		conf := patchConfigHeaders(nil, headers)
		if conf == nil {
			t.Error("empty request config")
			return
		}
		if conf.Headers["Accept"] != "application/json" {
			t.Error("incorrect headers")
		}
		conf2 := patchConfigHeaders(&RequestConfig{
			Headers: map[string]string{
				"Accept":       "application/xml",
				"Content-Type": "application/json",
			},
		}, headers)
		if conf2.Headers["Content-Type"] != "application/json" {
			t.Error("incorrect headers")
		}
		if conf2.Headers["Accept"] != "application/json" {
			t.Error("incorrect headers")
		}
	})

}
