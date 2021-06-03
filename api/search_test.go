package api

import (
	"bytes"
	"testing"
)

func TestSearch(t *testing.T) {
	checkClient(t)

	t.Run("Basic", func(t *testing.T) {
		sp := NewSP(spClient)
		data, err := sp.Search().PostQuery(&SearchQuery{
			QueryText: "*",
			RowLimit:  10,
		})
		if err != nil {
			t.Error(err)
		}
		if bytes.Compare(data, data.Normalized()) == -1 {
			t.Error("wrong response normalization")
		}
	})

	t.Run("Unmarshal", func(t *testing.T) {
		sp := NewSP(spClient)
		res, err := sp.Search().PostQuery(&SearchQuery{
			QueryText: "*",
			RowLimit:  10,
		})
		if err != nil {
			t.Error(err)
		}
		// if res.Data().ElapsedTime == 0 {
		// 	t.Error("incorrect response")
		// }
		if len(res.Data().PrimaryQueryResult.RelevantResults.Table.Rows) == 0 {
			t.Error("incorrect response")
		}
	})

	t.Run("Results", func(t *testing.T) {
		sp := NewSP(spClient)
		res, err := sp.Search().PostQuery(&SearchQuery{
			QueryText: "*",
			RowLimit:  10,
		})
		if err != nil {
			t.Error(err)
		}
		if len(res.Results()) == 0 {
			t.Error("incorrect response")
		}
	})

}
