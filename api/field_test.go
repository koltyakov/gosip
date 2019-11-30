package api

import (
	"testing"
)

func TestField(t *testing.T) {
	checkClient(t)

	web := NewSP(spClient).Web()
	field, err := getAnyField()
	if err != nil {
		t.Error(err)
	}

	t.Run("Get", func(t *testing.T) {
		if _, err := web.Fields().GetByID(field.ID).Get(); err != nil {
			t.Error(err)
		}
	})

}
