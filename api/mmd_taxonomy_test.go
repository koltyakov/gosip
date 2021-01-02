package api

import (
	"fmt"
	"strings"
	"testing"
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
}

func TestTaxonomyGroups(t *testing.T) {
	checkClient(t)

	taxonomy := NewSP(spClient).Taxonomy()

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
}

func TestTaxonomyTermSets(t *testing.T) {
	checkClient(t)

	taxonomy := NewSP(spClient).Taxonomy()

	termSetGUID, err := getTermSetID(taxonomy)
	if err != nil {
		t.Error(err)
	}

	t.Run("GetByID", func(t *testing.T) {
		data, err := taxonomy.Stores().Default().Sets().GetByID(termSetGUID).Select("Id").Get()
		if err != nil {
			t.Errorf("can't get term set by ID, %s\n", err)
		}
		if _, ok := data["Id"]; !ok {
			t.Error("should contain name property, cast error")
		}
	})

	t.Run("Terms/GetAll", func(t *testing.T) {
		_, err := taxonomy.Stores().Default().Sets().GetByID(termSetGUID).Terms().Select("Id,Name").GetAll()
		if err != nil {
			t.Errorf("%s", err)
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
	groupGUID, ok := gs[0]["Id"].(string)
	if !ok {
		return "", fmt.Errorf("can't get group ID")
	}
	return groupGUID, nil
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

	terms, err := taxonomy.Stores().Default().Sets().GetByID(termSetGUID).Terms().Select("Id").GetAll()
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
