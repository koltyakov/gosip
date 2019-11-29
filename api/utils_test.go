package api

import (
	"reflect"
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

}
