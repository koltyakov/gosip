package api

import (
	"fmt"
	"testing"
)

func TestOData(t *testing.T) {
	checkClient(t)

	t.Run("constructor", func(t *testing.T) {
		obj := &struct {
			modifiers *ODataMods
		}{
			modifiers: NewODataMods(),
		}
		obj.modifiers.AddSelect("Value")
		obj.modifiers.AddExpand("Value")
		if fmt.Sprintf("%+v", obj.modifiers.Get()) != "map[$expand:Value $select:Value]" {
			t.Error("incorrect add select result")
		}
	})

	t.Run("addSelect", func(t *testing.T) {
		modifiers := &ODataMods{}
		modifiers.AddSelect("Value")
		if fmt.Sprintf("%+v", modifiers.Get()) != "map[$select:Value]" {
			t.Error("incorrect add select result")
		}
	})

	t.Run("addExpand", func(t *testing.T) {
		modifiers := &ODataMods{}
		modifiers.AddExpand("Value")
		if fmt.Sprintf("%+v", modifiers.Get()) != "map[$expand:Value]" {
			t.Error("incorrect add expand result")
		}
	})

	t.Run("addFilter", func(t *testing.T) {
		modifiers := &ODataMods{}
		modifiers.AddFilter("Value")
		if fmt.Sprintf("%+v", modifiers.Get()) != "map[$filter:Value]" {
			t.Error("incorrect add filter result")
		}
	})

	t.Run("addSkip", func(t *testing.T) {
		modifiers := &ODataMods{}
		modifiers.AddSkip("Value")
		if fmt.Sprintf("%+v", modifiers.Get()) != "map[$skiptoken:Value]" {
			t.Error("incorrect add skip result")
		}
	})

	t.Run("addTop", func(t *testing.T) {
		modifiers := &ODataMods{}
		modifiers.AddTop(5)
		if fmt.Sprintf("%+v", modifiers.Get()) != "map[$top:5]" {
			t.Error("incorrect add top result")
		}
	})

	t.Run("addOrderBy", func(t *testing.T) {
		modifiers := &ODataMods{}
		modifiers.AddOrderBy("Field", true)
		if fmt.Sprintf("%+v", modifiers.Get()) != "map[$orderby:Field asc]" {
			t.Error("incorrect add order by result")
		}
	})

	t.Run("addMultyiOrderBy", func(t *testing.T) {
		modifiers := &ODataMods{}
		modifiers.AddOrderBy("OneField", true).AddOrderBy("AnotherField", false)
		if fmt.Sprintf("%+v", modifiers.Get()) != "map[$orderby:OneField asc,AnotherField desc]" {
			t.Error("incorrect add order by result")
		}
	})

	t.Run("addMixed", func(t *testing.T) {
		modifiers := &ODataMods{}
		modifiers.AddSelect("Select").AddExpand("Expand").AddTop(5)
		if fmt.Sprintf("%+v", modifiers.Get()) != "map[$expand:Expand $select:Select $top:5]" {
			t.Error("incorrect add mixed modifiers result")
		}
	})

	t.Run("toURL", func(t *testing.T) {
		modifiers := &ODataMods{}
		modifiers.AddSelect("Select").AddExpand("Expand").AddTop(5)
		if toURL("https://contoso/_api/Web", modifiers) != "https://contoso/_api/Web?%24expand=Expand&%24select=Select&%24top=5" {
			t.Error("incorrect add mixed modifiers result")
		}
	})

}
