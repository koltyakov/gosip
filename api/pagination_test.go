package api

import (
	"testing"
)

func TestPagination(t *testing.T) {
	checkClient(t)

	sp := NewSP(spClient)

	t.Run("Items", func(t *testing.T) {
		paged, err := sp.Web().GetList("_catalogs/masterpage").Items().Top(1).GetPaged()
		if err != nil {
			t.Fatal(err)
		}
		if paged.HasNextPage() {
			if _, err := paged.GetNextPage(); err != nil {
				t.Error(err)
			}
		} else {
			t.Logf("has no next page")
		}
	})

	// t.Run("Lists", func(t *testing.T) {
	// 	paged, err := sp.Web().Lists().Top(1).GetPaged()
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	// 	if paged.HasNextPage() {
	// 		if _, err := paged.GetNextPage(); err != nil {
	// 			t.Error(err)
	// 		}
	// 	} else {
	// 		t.Logf("has no next page")
	// 	}
	// })

	// t.Run("ContentTypes", func(t *testing.T) {
	// 	paged, err := sp.Web().ContentTypes().Top(1).GetPaged()
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	// 	if paged.HasNextPage() {
	// 		if _, err := paged.GetNextPage(); err != nil {
	// 			t.Error(err)
	// 		}
	// 	} else {
	// 		t.Logf("has no next page")
	// 	}
	// })

	// t.Run("Fields", func(t *testing.T) {
	// 	paged, err := sp.Web().Fields().Top(1).GetPaged()
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	// 	if paged.HasNextPage() {
	// 		if _, err := paged.GetNextPage(); err != nil {
	// 			t.Error(err)
	// 		}
	// 	} else {
	// 		t.Logf("has no next page")
	// 	}
	// })

	// t.Run("Files", func(t *testing.T) {
	// 	paged, err := sp.Web().GetFolder("_catalogs/masterpage").Files().Top(1).GetPaged()
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	// 	if paged.HasNextPage() {
	// 		if _, err := paged.GetNextPage(); err != nil {
	// 			t.Error(err)
	// 		}
	// 	} else {
	// 		t.Logf("has no next page")
	// 	}
	// })

	// t.Run("Folders", func(t *testing.T) {
	// 	paged, err := sp.Web().GetFolder("_catalogs/masterpage").Folders().Top(1).GetPaged()
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	// 	if paged.HasNextPage() {
	// 		if _, err := paged.GetNextPage(); err != nil {
	// 			t.Error(err)
	// 		}
	// 	} else {
	// 		t.Logf("has no next page")
	// 	}
	// })

	// t.Run("Groups", func(t *testing.T) {
	// 	paged, err := sp.Web().SiteGroups().Top(1).GetPaged()
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	// 	if paged.HasNextPage() {
	// 		if _, err := paged.GetNextPage(); err != nil {
	// 			t.Error(err)
	// 		}
	// 	} else {
	// 		t.Logf("has no next page")
	// 	}
	// })

	// t.Run("RecycleBin", func(t *testing.T) {
	// 	paged, err := sp.Web().RecycleBin().Top(1).GetPaged()
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	// 	if paged.HasNextPage() {
	// 		if _, err := paged.GetNextPage(); err != nil {
	// 			t.Error(err)
	// 		}
	// 	} else {
	// 		t.Logf("has no next page")
	// 	}
	// })

	// t.Run("Users", func(t *testing.T) {
	// 	paged, err := sp.Web().SiteUsers().Top(1).GetPaged()
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	// 	if paged.HasNextPage() {
	// 		if _, err := paged.GetNextPage(); err != nil {
	// 			t.Error(err)
	// 		}
	// 	} else {
	// 		t.Logf("has no next page")
	// 	}
	// })

	// t.Run("Views", func(t *testing.T) {
	// 	paged, err := sp.Web().GetList("_catalogs/masterpage").Views().Top(1).GetPaged()
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	// 	if paged.HasNextPage() {
	// 		if _, err := paged.GetNextPage(); err != nil {
	// 			t.Error(err)
	// 		}
	// 	} else {
	// 		t.Logf("has no next page")
	// 	}
	// })

	// t.Run("Webs", func(t *testing.T) {
	// 	paged, err := sp.Web().Webs().Top(1).GetPaged()
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	// 	if paged.HasNextPage() {
	// 		if _, err := paged.GetNextPage(); err != nil {
	// 			t.Error(err)
	// 		}
	// 	} else {
	// 		t.Logf("has no next page")
	// 	}
	// })

}
