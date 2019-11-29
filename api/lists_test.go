package api

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestLists(t *testing.T) {
	checkClient(t)

	web := NewSP(spClient).Web()
	newListTitle := uuid.New().String()

	t.Run("Get", func(t *testing.T) {
		data, err := web.Lists().Select("Id,Title").Conf(headers.verbose).Get()
		if err != nil {
			t.Error(err)
		}

		res := &struct {
			D struct {
				Results []interface{} `json:"results"`
			} `json:"d"`
		}{}

		if err := json.Unmarshal(data, &res); err != nil {
			t.Error(err)
		}

		if len(res.D.Results) == 0 {
			t.Error("can't get webs")
		}
	})

	t.Run("Add", func(t *testing.T) {
		if _, err := web.Lists().Add(newListTitle, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		listID, _, err := getAnyList()
		if err != nil {
			t.Error(err)
		}
		if _, err := web.Lists().GetByID(listID).Get(); err != nil {
			t.Error(err)
		}
	})

	t.Run("GetByTitle", func(t *testing.T) {
		_, title, err := getAnyList()
		if err != nil {
			t.Error(err)
		}
		if _, err := web.Lists().GetByTitle(title).Get(); err != nil {
			t.Error(err)
		}
	})

	t.Run("GetListByURI", func(t *testing.T) {
		listURI := getRelativeURL(spClient.AuthCnfg.GetSiteURL()) +
			"/Lists/" + strings.ReplaceAll(newListTitle, "-", "")
		if _, err := web.GetList(listURI).Get(); err != nil {
			t.Error(err)
		}
	})

	t.Run("AddWithURI", func(t *testing.T) {
		listTitle := uuid.New().String()
		listURI := uuid.New().String()
		if _, err := web.Lists().AddWithURI(listTitle, listURI, nil); err != nil {
			t.Error(err)
		}
		if _, err := web.Lists().GetByTitle(listTitle).Delete(); err != nil {
			t.Error(err)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		if _, err := web.Lists().GetByTitle(newListTitle).Delete(); err != nil {
			t.Error(err)
		}
	})

}

func getAnyList() (string, string, error) {
	web := NewSP(spClient).Web()
	data, err := web.Lists().Select("Id,Title").Top(1).Conf(headers.verbose).Get()
	if err != nil {
		return "", "", err
	}

	res := &struct {
		D struct {
			Results []struct {
				ID    string `json:"Id"`
				Title string `json:"Title"`
			} `json:"results"`
		} `json:"d"`
	}{}

	if err := json.Unmarshal(data, &res); err != nil {
		return "", "", err
	}

	if len(res.D.Results) == 0 {
		return "", "", fmt.Errorf("can't get webs")
	}

	return res.D.Results[0].ID, res.D.Results[0].Title, nil
}
