package api

import (
	"testing"
)

func TestContentType(t *testing.T) {
	checkClient(t)

	web := NewSP(spClient).Web()

	t.Run("Get", func(t *testing.T) {
		if _, err := web.ContentTypes().Conf(headers.verbose).Top(1).Get(); err != nil {
			t.Error(err)
		}
	})

}
