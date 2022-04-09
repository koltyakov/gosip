package api

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestRecordsModern(t *testing.T) {
	checkClient(t)

	if envCode != "spo" {
		t.Skip("is not supported")
	}

	sp := NewSP(spClient)

	// Activate in place record management feature
	if err := sp.Site().Features().Add("da2e115b-07e4-49d9-bb2c-35e93bb9fca9", true); err != nil {
		t.Error(err)
	}

	folder := sp.Web().GetFolder("Shared Documents")

	folderName := uuid.New().String()
	docs := []string{
		fmt.Sprintf("%s.txt", uuid.New().String()),
	}

	if _, err := folder.Folders().Add(folderName); err != nil {
		t.Error(err)
	}

	for _, doc := range docs {
		if _, err := folder.Folders().GetByName(folderName).Files().Add(doc, []byte(doc), true); err != nil {
			t.Error(err)
		}
	}

	t.Run("Records/LockRecordItem", func(t *testing.T) {
		item, err := folder.Folders().GetByName(folderName).Files().GetByName(docs[0]).GetItem()
		if err != nil {
			t.Error(err)
		}
		if err := item.Records().LockRecordItem(); err != nil {
			t.Error(err)
		}
	})

	t.Run("Records/UnlockRecordItem", func(t *testing.T) {
		item, err := folder.Folders().GetByName(folderName).Files().GetByName(docs[0]).GetItem()
		if err != nil {
			t.Error(err)
		}
		if err := item.Records().UnlockRecordItem(); err != nil {
			t.Error(err)
		}
	})

	if err := folder.Folders().GetByName(folderName).Delete(); err != nil {
		t.Error(err)
	}

}
