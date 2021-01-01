package api

import (
	"testing"
)

func TestTaxonomy(t *testing.T) {
	checkClient(t)

	taxonomy := NewSP(spClient).Taxonomy()

	t.Run("TermStore", func(t *testing.T) {
		ts := taxonomy.TermStore()
		if ts == nil {
			t.Error("can't get term store object")
		}
		tsInfo, err := ts.Get()
		if err != nil {
			t.Errorf("can't get term store info, %s\n", err)
		}
		if _, ok := tsInfo["Name"]; !ok {
			t.Error("should contain name property, cast error")
		}

		t.Run("ByID", func(t *testing.T) {
			data, err := taxonomy.TermStoreByID(tsInfo["Id"].(string)).Get()
			if err != nil {
				t.Errorf("can't get term store info, %s\n", err)
			}
			if _, ok := data["Id"]; !ok {
				t.Error("should contain name property, cast error")
			}
		})

		t.Run("ByName", func(t *testing.T) {
			data, err := taxonomy.TermStoreByName(tsInfo["Name"].(string)).Get()
			if err != nil {
				t.Errorf("can't get term store info, %s\n", err)
			}
			if _, ok := data["Id"]; !ok {
				t.Error("should contain name property, cast error")
			}
		})

		t.Run("Select", func(t *testing.T) {
			ts := taxonomy.TermStore().Select("Id,Name,Groups")
			if len(ts.selectProps) != 3 {
				t.Error("error setting props")
			}
			ts.Select("IsOnline")
			if len(ts.selectProps) != 4 {
				t.Error("error setting props")
			}
			ts.Select("Id")
			if len(ts.selectProps) != 4 {
				t.Error("error setting props")
			}
			data, err := ts.Get()
			if err != nil {
				t.Errorf("can't get term store info, %s\n", err)
			}
			if _, ok := data["Name"]; !ok {
				t.Error("should contain name property, cast error")
			}
			if _, ok := data["Groups"]; !ok {
				t.Error("should contain name property, cast error")
			}
		})
	})

	t.Run("TermGroups", func(t *testing.T) {
		gs, err := taxonomy.TermStore().Groups().Get()
		if err != nil {
			t.Error(err)
		}
		if len(gs) == 0 {
			t.Error("can't get term store groups")
		}

		groupGUID, ok := gs[0]["Id"].(string)
		if !ok {
			t.Error("can't get group ID")
		}

		group, err := taxonomy.TermStore().Groups().GetByID(groupGUID).Select("Id,Name").Get()
		if err != nil {
			t.Error(err)
		}

		if group["Id"].(string) != groupGUID {
			t.Error("error getting group info")
		}
	})

	t.Run("TermSets", func(t *testing.T) {
		gs, err := taxonomy.TermStore().Groups().Get()
		if err != nil {
			t.Error(err)
		}
		if len(gs) == 0 {
			t.Error("can't get term store groups")
		}

		groupGUID, ok := gs[0]["Id"].(string)
		if !ok {
			t.Error("can't get group ID")
		}

		termSets := taxonomy.TermStore().Groups().GetByID(groupGUID).TermSets()
		if termSets == nil {
			t.Error("can't get term sets for a group")
		}

		tsData, err := termSets.Get()
		if err != nil {
			t.Error(err)
		}
		if len(tsData) == 0 {
			t.Error("can't get group term sets")
		}

		ts := tsData[0]

		t.Run("ByID", func(t *testing.T) {
			data, err := taxonomy.TermStore().GetTermSet(ts["Id"].(string)).Get()
			if err != nil {
				t.Errorf("can't get term set by ID, %s\n", err)
			}
			if _, ok := data["Id"]; !ok {
				t.Error("should contain name property, cast error")
			}
		})

		t.Run("Terms/GetAll", func(t *testing.T) {
			_, err := taxonomy.TermStore().GetTermSet(ts["Id"].(string)).Terms().GetAll()
			if err != nil {
				t.Errorf("%s", err)
			}
		})
	})

}

func TestTaxonomyUtils(t *testing.T) {
	t.Run("Taxonomy/AppendProp", func(t *testing.T) {
		props := appendProp([]string{}, "Id,Name,Groups")
		if len(props) != 3 {
			t.Error("error setting props")
		}
		props = appendProp(props, "IsOnline")
		if len(props) != 4 {
			t.Error("error setting props")
		}
		props = appendProp(props, "Id")
		if len(props) != 4 {
			t.Error("error setting props")
		}
	})

	t.Run("Taxonomy/TrimGuid", func(t *testing.T) {
		if trimGUID("/Guid(9dd47937-e620-4196-87a7-815c7e6aa384)/") != "9dd47937-e620-4196-87a7-815c7e6aa384" {
			t.Error("error trimming GUID")
		}
		if trimGUID("/guid(9dd47937-e620-4196-87a7-815c7e6aa384)/") != "9dd47937-e620-4196-87a7-815c7e6aa384" {
			t.Error("error trimming GUID")
		}
		if trimGUID("/GUID(9dd47937-e620-4196-87a7-815c7e6aa384)/") != "9dd47937-e620-4196-87a7-815c7e6aa384" {
			t.Error("error trimming GUID")
		}
	})
}
