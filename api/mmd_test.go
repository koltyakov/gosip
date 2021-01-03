package api

import (
	"fmt"
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
}

func TestTaxonomyGroups(t *testing.T) {
	checkClient(t)

	taxonomy := NewSP(spClient).Taxonomy()

	newGroupGUID := uuid.New().String()
	newGroupName := "Delete me " + newGroupGUID

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
		groupGUID, err := getTermGroupID(taxonomy)
		if err != nil {
			t.Error(err)
		}

		group, err := taxonomy.Stores().Default().Groups().GetByID(groupGUID).Select("Id,Name").Get()
		if err != nil {
			t.Error(err)
		}

		if group["Id"].(string) != groupGUID {
			t.Error("error getting group info")
		}
	})

	t.Run("Add", func(t *testing.T) {
		group, err := taxonomy.Stores().Default().Groups().Add(newGroupName, newGroupGUID)
		if err != nil {
			t.Error(err)
		}
		if group["Name"].(string) != newGroupName {
			t.Error("error getting group info")
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

	termGroupID, err := getTermGroupID(taxonomy)
	if err != nil {
		t.Error(err)
	}

	termSetGUID, err := getTermSetID(taxonomy)
	if err != nil {
		t.Error(err)
	}

	newTermSetGUID := uuid.New().String()
	newTermSetName := "Delete me " + newTermSetGUID

	t.Run("GetByID", func(t *testing.T) {
		data, err := taxonomy.Stores().Default().Sets().GetByID(termSetGUID).Select("Id").Get()
		if err != nil {
			t.Errorf("can't get term set by ID, %s\n", err)
		}
		if _, ok := data["Id"]; !ok {
			t.Error("should contain name property, cast error")
		}
	})

	t.Run("Terms/GetAllTerms", func(t *testing.T) {
		_, err := taxonomy.Stores().Default().Sets().GetByID(termSetGUID).Select("Id,Name").GetAllTerms()
		if err != nil {
			t.Errorf("%s", err)
		}
	})

	t.Run("Terms/Get", func(t *testing.T) {
		_, err := taxonomy.Stores().Default().Sets().GetByID(termSetGUID).Terms().Select("Id,Name").Get()
		if err != nil {
			t.Errorf("%s", err)
		}
	})

	t.Run("Add", func(t *testing.T) {
		store := taxonomy.Stores().Default()

		tsInfo, err := store.Select("DefaultLanguage").Get()
		if err != nil {
			t.Error(err)
		}
		lang := int(tsInfo["DefaultLanguage"].(float64))

		termSet, err := store.Groups().GetByID(termGroupID).Sets().Add(newTermSetName, newTermSetGUID, lang)
		if err != nil {
			t.Error(err)
		}
		if termSet["Name"].(string) != newTermSetName {
			t.Error("error getting term set info")
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

	termSetGUID, err := getTermSetID(taxonomy)
	if err != nil {
		t.Error(err)
	}

	termGUID, err := getTermID(taxonomy)
	if err != nil {
		t.Error(err)
	}

	t.Run("Store/GetTerm", func(t *testing.T) {
		termInfo, err := taxonomy.Stores().Default().Terms().GetByID(termGUID).Select("Id").Get()
		if err != nil {
			t.Error(err)
		}
		if termGUID != termInfo["Id"].(string) {
			t.Error("unexpected term ID")
		}
	})

	t.Run("TermSet/GetByID", func(t *testing.T) {
		termInfo, err := taxonomy.Stores().Default().Sets().GetByID(termSetGUID).Terms().GetByID(termGUID).Select("Id").Get()
		if err != nil {
			t.Error(err)
		}
		if termGUID != termInfo["Id"].(string) {
			t.Error("unexpected term ID")
		}
	})

	t.Run("Terms/CRUD", func(t *testing.T) {
		tsInfo, err := taxonomy.Stores().Default().Select("DefaultLanguage").Get()
		if err != nil {
			t.Error(err)
		}
		lang := int(tsInfo["DefaultLanguage"].(float64))

		newTermGUID := uuid.New().String()
		newTermName := "Delete me " + newTermGUID

		t.Run("Add", func(t *testing.T) {
			termSet := taxonomy.Stores().Default().Sets().GetByID(termSetGUID)
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

			subTerms, err := parentTerm.Terms().Get() // All()
			if err != nil {
				t.Error(err)
			}
			if len(subTerms) != 1 {
				t.Error("error getting subterms")
			}
		})

		t.Run("Add#FailAddingADuplicate", func(t *testing.T) {
			termSet := taxonomy.Stores().Default().Sets().GetByID(termSetGUID)
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
		if err != nil && strings.Index(err.Error(), "Guid should contain 32 digits with 4 dashes") == -1 {
			t.Error("should fail with Guid should contain 32 digits with 4 dashes message")
		}

		_, err = store.Groups().GetByID("c5589d9f-8044-b000-5f6d-bcc9f93b8b8d").Get()
		if err == nil {
			t.Error("should fail with object not found message")
		}
		if err != nil && strings.Index(err.Error(), "object not found") == -1 {
			t.Error("should fail with object not found message")
		}
	})

	t.Run("Sets/GetByID", func(t *testing.T) {
		_, err := store.Sets().GetByID("wrong-id").Get()
		if err == nil {
			t.Error("should fail with Guid should contain 32 digits with 4 dashes message")
		}
		if err != nil && strings.Index(err.Error(), "Guid should contain 32 digits with 4 dashes") == -1 {
			t.Error("should fail with Guid should contain 32 digits with 4 dashes message")
		}

		_, err = store.Sets().GetByID("c5589d9f-8044-b000-5f6d-bcc9f93b8b8d").Get()
		if err == nil {
			t.Error("should fail with object not found message")
		}
		if err != nil && strings.Index(err.Error(), "object not found") == -1 {
			t.Error("should fail with object not found message")
		}
	})

	t.Run("Terms/GetByID", func(t *testing.T) {
		_, err := store.Terms().GetByID("wrong-id").Get()
		if err == nil {
			t.Error("should fail with Guid should contain 32 digits with 4 dashes message")
		}
		if err != nil && strings.Index(err.Error(), "Guid should contain 32 digits with 4 dashes") == -1 {
			t.Error("should fail with Guid should contain 32 digits with 4 dashes message")
		}

		_, err = store.Terms().GetByID("c5589d9f-8044-b000-5f6d-bcc9f93b8b8d").Get()
		if err == nil {
			t.Error("should fail with object not found message")
		}
		if err != nil && strings.Index(err.Error(), "object not found") == -1 {
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

func getTermGroupID(taxonomy *Taxonomy) (string, error) {
	gs, err := taxonomy.Stores().Default().Groups().Get()
	if err != nil {
		return "", err
	}
	if len(gs) == 0 {
		return "", fmt.Errorf("can't get term store groups")
	}
	for _, group := range gs {
		if strings.Index(group["Name"].(string), "Delete me ") == -1 {
			groupGUID, ok := group["Id"].(string)
			if !ok {
				return "", fmt.Errorf("can't get group ID")
			}
			return groupGUID, nil
		}
	}
	return "", fmt.Errorf("can't get group ID")
}

func getTermSetID(taxonomy *Taxonomy) (string, error) {
	groupGUID, err := getTermGroupID(taxonomy)
	if err != nil {
		return "", err
	}

	termSets := taxonomy.Stores().Default().Groups().GetByID(groupGUID).Sets()
	if termSets == nil {
		return "", err
	}

	tsData, err := termSets.Get()
	if err != nil {
		return "", err
	}
	if len(tsData) == 0 {
		return "", fmt.Errorf("can't get term sets")
	}

	termSetGUID, ok := tsData[0]["Id"].(string)
	if !ok {
		return "", fmt.Errorf("can't get term set ID")
	}

	return termSetGUID, nil
}

func getTermID(taxonomy *Taxonomy) (string, error) {
	termSetGUID, err := getTermSetID(taxonomy)
	if err != nil {
		return "", err
	}

	terms, err := taxonomy.Stores().Default().Sets().GetByID(termSetGUID).Select("Id").GetAllTerms()
	if terms == nil {
		return "", err
	}
	if len(terms) == 0 {
		return "", fmt.Errorf("can't get term sets")
	}

	termGUID, ok := terms[0]["Id"].(string)
	if !ok {
		return "", fmt.Errorf("can't get term ID")
	}

	return termGUID, nil
}
