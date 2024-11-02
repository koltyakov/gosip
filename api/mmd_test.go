package api

import (
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestTaxonomyStores(t *testing.T) {
	checkClient(t)

	taxonomy := NewSP(spClient).Taxonomy()

	t.Run("Default", func(t *testing.T) {
		tsInfo, err := taxonomy.Stores().Default().Get()
		if err != nil {
			t.Errorf("can't get term store info, %s\n", err)
		}
		if _, ok := tsInfo["Name"]; !ok {
			t.Error("should contain name property, cast error")
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		tsInfo, err := taxonomy.Stores().Default().Get()
		if err != nil {
			t.Errorf("can't get term store info, %s\n", err)
		}
		data, err := taxonomy.Stores().GetByID(tsInfo["Id"].(string)).Get()
		if err != nil {
			t.Errorf("can't get term store info, %s\n", err)
		}
		if _, ok := data["Id"]; !ok {
			t.Error("should contain name property, cast error")
		}
	})

	t.Run("GetByName", func(t *testing.T) {
		tsInfo, err := taxonomy.Stores().Default().Get()
		if err != nil {
			t.Errorf("can't get term store info, %s\n", err)
		}
		data, err := taxonomy.Stores().GetByName(tsInfo["Name"].(string)).Get()
		if err != nil {
			t.Errorf("can't get term store info, %s\n", err)
		}
		if _, ok := data["Id"]; !ok {
			t.Error("should contain name property, cast error")
		}
	})

	t.Run("Select", func(t *testing.T) {
		ts := taxonomy.Stores().Default().Select("Id,Name,Groups")
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

	t.Run("UpdateCache", func(t *testing.T) {
		if err := taxonomy.Stores().Default().UpdateCache(); err != nil {
			t.Errorf("can't get term store info, %s\n", err)
		}
	})

	t.Run("Sets/GetByName", func(t *testing.T) {
		sets, err := taxonomy.Stores().Default().Sets().GetByName("Department", 1033)
		if err != nil {
			t.Errorf("can't get term set by name, %s\n", err)
		}
		if len(sets) != 1 {
			t.Log("maybe can't get term set by name")
		}
	})
}

func TestTaxonomyGroups(t *testing.T) {
	checkClient(t)

	taxonomy := NewSP(spClient).Taxonomy()

	newGroupGUID := uuid.New().String()
	newGroupName := "Delete me " + newGroupGUID

	t.Run("Add", func(t *testing.T) {
		group, err := taxonomy.Stores().Default().Groups().Add(newGroupName, newGroupGUID)
		if err != nil {
			t.Error(err)
		}
		if group["Name"].(string) != newGroupName {
			t.Error("error getting group info")
		}
	})

	t.Run("Get", func(t *testing.T) {
		gs, err := taxonomy.Stores().Default().Groups().Get()
		if err != nil {
			t.Error(err)
		}
		if len(gs) == 0 {
			t.Error("can't get term store groups")
		}
		if _, ok := gs[0]["Id"].(string); !ok {
			t.Error("can't get group ID")
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		group, err := taxonomy.Stores().Default().Groups().GetByID(newGroupGUID).Select("Id,Name").Get()
		if err != nil {
			t.Error(err)
		}

		if newGroupGUID != trimTaxonomyGUID(group["Id"].(string)) {
			t.Error("error getting group info")
		}
	})

	t.Run("Sets/GetByName", func(t *testing.T) {
		sets, err := taxonomy.Stores().Default().Groups().GetByID("any-id-should-work-here").Sets().GetByName("Department", 1033)
		if err != nil {
			t.Errorf("can't get term set by name, %s\n", err)
		}
		if len(sets) != 1 {
			t.Log("maybe can't get term set by name")
		}
	})

	t.Run("Delete", func(t *testing.T) {
		if err := taxonomy.Stores().Default().Groups().GetByID(newGroupGUID).Delete(); err != nil {
			t.Error(err)
		}
	})
}

func TestTaxonomyTermSets(t *testing.T) {
	checkClient(t)

	taxonomy := NewSP(spClient).Taxonomy()

	lang, err := getDefaultLang(taxonomy)
	if err != nil {
		t.Error(err)
	}

	newGroupGUID := uuid.New().String()
	newGroupName := "Delete me " + newGroupGUID

	_, err = taxonomy.Stores().Default().Groups().Add(newGroupName, newGroupGUID)
	if err != nil {
		t.Error(err)
	}

	newTermSetGUID := uuid.New().String()
	newTermSetName := "Delete me " + newTermSetGUID

	defer func() {
		if err := taxonomy.Stores().Default().Groups().GetByID(newGroupGUID).Delete(); err != nil {
			t.Error(err)
		}
	}()

	t.Run("Add", func(t *testing.T) {
		store := taxonomy.Stores().Default()
		termSet, err := store.Groups().GetByID(newGroupGUID).Sets().Add(newTermSetName, newTermSetGUID, lang)
		if err != nil {
			t.Error(err)
		}
		if termSet["Name"].(string) != newTermSetName {
			t.Error("error getting term set info")
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		data, err := taxonomy.Stores().Default().Sets().GetByID(newTermSetGUID).Select("Id").Get()
		if err != nil {
			t.Errorf("can't get term set by ID, %s\n", err)
		}
		if _, ok := data["Id"]; !ok {
			t.Error("should contain name property, cast error")
		}
	})

	t.Run("GetAllTerms", func(t *testing.T) {
		_, err := taxonomy.Stores().Default().Sets().GetByID(newTermSetGUID).Select("Id,Name").GetAllTerms()
		if err != nil {
			t.Errorf("%s", err)
		}
	})

	t.Run("Terms/Get", func(t *testing.T) {
		_, err := taxonomy.Stores().Default().Sets().GetByID(newTermSetGUID).Terms().Select("Id,Name").Get()
		if err != nil {
			t.Errorf("%s", err)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		store := taxonomy.Stores().Default()
		if err := store.Sets().GetByID(newTermSetGUID).Delete(); err != nil {
			t.Error(err)
		}
	})
}

func TestTaxonomyTerms(t *testing.T) {
	checkClient(t)

	taxonomy := NewSP(spClient).Taxonomy()

	lang, err := getDefaultLang(taxonomy)
	if err != nil {
		t.Error(err)
	}

	newGroupGUID := uuid.New().String()
	newGroupName := "Delete me " + newGroupGUID

	_, err = taxonomy.Stores().Default().Groups().Add(newGroupName, newGroupGUID)
	if err != nil {
		t.Error(err)
	}

	newTermSetGUID := uuid.New().String()
	newTermSetName := "Delete me " + newTermSetGUID

	_, err = taxonomy.Stores().Default().Groups().GetByID(newGroupGUID).Sets().Add(newTermSetName, newTermSetGUID, lang)
	if err != nil {
		t.Error(err)
	}

	newTermGUID := uuid.New().String()
	newTermName := "Delete me " + newTermSetGUID

	_, err = taxonomy.Stores().Default().Sets().GetByID(newTermSetGUID).Terms().Add(newTermName, newTermGUID, lang)
	if err != nil {
		t.Error(err)
	}

	defer func() {
		if err := taxonomy.Stores().Default().Groups().GetByID(newGroupGUID).Delete(); err != nil {
			t.Error(err)
		}
	}()

	t.Run("Store/GetTerm", func(t *testing.T) {
		termInfo, err := taxonomy.Stores().Default().Terms().GetByID(newTermGUID).Select("Id").Get()
		if err != nil {
			t.Error(err)
		}
		if newTermGUID != trimTaxonomyGUID(termInfo["Id"].(string)) {
			t.Error("unexpected term ID")
		}
	})

	t.Run("TermSet/GetByID", func(t *testing.T) {
		termSet := taxonomy.Stores().Default().Sets().GetByID(newTermSetGUID)
		termInfo, err := termSet.Terms().GetByID(newTermGUID).Select("Id").Get()
		if err != nil {
			t.Error(err)
		}
		if newTermGUID != trimTaxonomyGUID(termInfo["Id"].(string)) {
			t.Error("unexpected term ID")
		}
	})

	t.Run("Terms/CRUD", func(t *testing.T) {
		newTermGUID := uuid.New().String()
		newTermName := "Delete me " + newTermGUID

		t.Run("Add", func(t *testing.T) {
			termSet := taxonomy.Stores().Default().Sets().GetByID(newTermSetGUID)
			termInfo, err := termSet.Terms().Add(newTermName, newTermGUID, lang)
			if err != nil {
				t.Error(err)
			}
			if newTermGUID != trimTaxonomyGUID(termInfo["Id"].(string)) {
				t.Error("unexpected term ID")
			}
		})

		t.Run("Add#ChildTerm", func(t *testing.T) {
			subTermGUID := uuid.New().String()
			subTermName := "Sub term " + subTermGUID

			store := taxonomy.Stores().Default()
			parentTerm := store.Terms().GetByID(newTermGUID)

			termInfo, err := parentTerm.Terms().Add(subTermName, subTermGUID, lang)
			if err != nil {
				t.Error(err)
			}
			if subTermGUID != trimTaxonomyGUID(termInfo["Id"].(string)) {
				t.Error("unexpected term ID")
			}

			subTerms, err := parentTerm.Terms().Get()
			if err != nil {
				t.Error(err)
			}
			if len(subTerms) != 1 {
				t.Log("error getting subterms")
			}
		})

		t.Run("Move/ToTerm", func(t *testing.T) {
			childTermGUID := uuid.New().String()
			childTermName := "Movable term " + childTermGUID

			termSet := taxonomy.Stores().Default().Sets().GetByID(newTermSetGUID)

			if _, err := termSet.Terms().Add(childTermName, childTermGUID, lang); err != nil {
				t.Error(err)
			}

			childTerm := termSet.Terms().GetByID(childTermGUID)
			if err := childTerm.Move(newTermSetGUID, newTermGUID); err != nil {
				t.Error(err)
			}
		})

		t.Run("Move/ToTermSet", func(t *testing.T) {
			childTermGUID := uuid.New().String()
			childTermName := "Movable term " + childTermGUID

			termSet := taxonomy.Stores().Default().Sets().GetByID(newTermSetGUID)

			if _, err := termSet.Terms().GetByID(newTermGUID).Terms().Add(childTermName, childTermGUID, lang); err != nil {
				t.Error(err)
			}

			childTerm := termSet.Terms().GetByID(childTermGUID)
			if err := childTerm.Move(newTermSetGUID, ""); err != nil {
				t.Error(err)
			}

			if err := childTerm.Delete(); err != nil {
				t.Error(err)
			}
		})

		t.Run("Add#FailAddingADuplicate", func(t *testing.T) {
			termSet := taxonomy.Stores().Default().Sets().GetByID(newTermSetGUID)
			if _, err := termSet.Terms().Add(newTermName, newTermGUID, lang); err == nil {
				t.Error("should fail with duplicate error message")
			}
		})

		t.Run("Get", func(t *testing.T) {
			term := taxonomy.Stores().Default().Terms().GetByID(newTermGUID)
			if _, err := term.Get(); err != nil {
				t.Error(err)
			}
		})

		t.Run("Update", func(t *testing.T) {
			updateTermName := newTermName + " (updated)"
			props := map[string]interface{}{"Name": updateTermName}
			term := taxonomy.Stores().Default().Terms().GetByID(newTermGUID)
			termAfterUpdate, err := term.Update(props)
			if err != nil {
				t.Error(err)
			}
			if termAfterUpdate["Name"].(string) != updateTermName {
				t.Error("failed to update term name")
			}
		})

		t.Run("Deprecate", func(t *testing.T) {
			store := taxonomy.Stores().Default()
			term := store.Terms().GetByID(newTermGUID)
			if err := term.Deprecate(true); err != nil {
				t.Error(err)
			}
			// if err := store.UpdateCache(); err != nil {
			// 	t.Error(err)
			// }
			// data, err := term.Select("IsDeprecated").Get()
			// if err != nil {
			// 	t.Error(err)
			// }
			// if !data["IsDeprecated"].(bool) {
			// 	// t.Error("failed to deprecate")
			// 	t.Log("maybe failed to deprecate")
			// }
		})

		t.Run("Activate", func(t *testing.T) {
			store := taxonomy.Stores().Default()
			term := store.Terms().GetByID(newTermGUID)
			if err := term.Deprecate(false); err != nil {
				t.Error(err)
			}
			// if err := store.UpdateCache(); err != nil {
			// 	t.Error(err)
			// }
			// data, err := term.Select("IsDeprecated").Get()
			// if err != nil {
			// 	t.Error(err)
			// }
			// if data["IsDeprecated"].(bool) {
			// 	// t.Error("failed to activate")
			// 	t.Log("maybe failed to activate")
			// }
		})

		t.Run("Delete", func(t *testing.T) {
			term := taxonomy.Stores().Default().Terms().GetByID(newTermGUID)
			if err := term.Delete(); err != nil {
				t.Error(err)
			}
			if err := taxonomy.Stores().Default().Sets().GetByID(newTermSetGUID).Delete(); err != nil {
				t.Error(err)
			}
		})
	})
}

func TestTaxonomyNotFoundIDs(t *testing.T) {
	checkClient(t)

	taxonomy := NewSP(spClient).Taxonomy()
	store := taxonomy.Stores().Default()

	t.Run("Groups/GetByID", func(t *testing.T) {
		_, err := store.Groups().GetByID("wrong-id").Get()
		if err == nil {
			t.Error("should fail with Guid should contain 32 digits with 4 dashes message")
		}
		if err != nil && !strings.Contains(err.Error(), "Guid should contain 32 digits with 4 dashes") {
			t.Error("should fail with Guid should contain 32 digits with 4 dashes message")
		}

		_, err = store.Groups().GetByID("c5589d9f-8044-b000-5f6d-bcc9f93b8b8d").Get()
		if err == nil {
			t.Error("should fail with object not found message")
		}
		if err != nil && !strings.Contains(err.Error(), "object not found") {
			t.Error("should fail with object not found message")
		}
	})

	t.Run("Sets/GetByID", func(t *testing.T) {
		_, err := store.Sets().GetByID("wrong-id").Get()
		if err == nil {
			t.Error("should fail with Guid should contain 32 digits with 4 dashes message")
		}
		if err != nil && !strings.Contains(err.Error(), "Guid should contain 32 digits with 4 dashes") {
			t.Error("should fail with Guid should contain 32 digits with 4 dashes message")
		}

		_, err = store.Sets().GetByID("c5589d9f-8044-b000-5f6d-bcc9f93b8b8d").Get()
		if err == nil {
			t.Error("should fail with object not found message")
		}
		if err != nil && !strings.Contains(err.Error(), "object not found") {
			t.Error("should fail with object not found message")
		}
	})

	t.Run("Terms/GetByID", func(t *testing.T) {
		_, err := store.Terms().GetByID("wrong-id").Get()
		if err == nil {
			t.Error("should fail with Guid should contain 32 digits with 4 dashes message")
		}
		if err != nil && !strings.Contains(err.Error(), "Guid should contain 32 digits with 4 dashes") {
			t.Error("should fail with Guid should contain 32 digits with 4 dashes message")
		}

		_, err = store.Terms().GetByID("c5589d9f-8044-b000-5f6d-bcc9f93b8b8d").Get()
		if err == nil {
			t.Error("should fail with object not found message")
		}
		if err != nil && !strings.Contains(err.Error(), "object not found") {
			t.Error("should fail with object not found message")
		}
	})
}

func TestTaxonomyUtils(t *testing.T) {
	t.Run("Taxonomy/AppendProp", func(t *testing.T) {
		props := appendTaxonomyProp([]string{}, "Id,Name,Groups")
		if len(props) != 3 {
			t.Error("error setting props")
		}
		props = appendTaxonomyProp(props, "IsOnline")
		if len(props) != 4 {
			t.Error("error setting props")
		}
		props = appendTaxonomyProp(props, "Id")
		if len(props) != 4 {
			t.Error("error setting props")
		}
	})

	t.Run("Taxonomy/TrimGuid", func(t *testing.T) {
		if trimTaxonomyGUID("/Guid(9dd47937-e620-4196-87a7-815c7e6aa384)/") != "9dd47937-e620-4196-87a7-815c7e6aa384" {
			t.Error("error trimming GUID")
		}
		if trimTaxonomyGUID("/guid(9dd47937-e620-4196-87a7-815c7e6aa384)/") != "9dd47937-e620-4196-87a7-815c7e6aa384" {
			t.Error("error trimming GUID")
		}
		if trimTaxonomyGUID("/GUID(9dd47937-e620-4196-87a7-815c7e6aa384)/") != "9dd47937-e620-4196-87a7-815c7e6aa384" {
			t.Error("error trimming GUID")
		}
	})
}

func getDefaultLang(taxonomy *Taxonomy) (int, error) {
	tsInfo, err := taxonomy.Stores().Default().Select("DefaultLanguage").Get()
	if err != nil {
		return 1033, err
	}
	lang := int(tsInfo["DefaultLanguage"].(float64))
	return lang, nil
}
