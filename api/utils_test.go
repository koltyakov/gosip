package api

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
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
		resStr := trimMultiline(initial)
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

		// should not add __metadata if present
		m2 := []byte(`{"__metadata":{"type":"SP.Any"},"prop":"val"}`)
		p2 := patchMetadataType(m2, "SP.List")
		if !containsMetadataType(p2) {
			t.Error("payload lost metadata prop")
		}
		if strings.Index(fmt.Sprintf("%s", p2), `"__metadata":{"type":"SP.Any"}`) == -1 {
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

}
