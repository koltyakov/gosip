package handlers

import (
	"fmt"
	"net/http"

	"github.com/koltyakov/gosip"
)

// GetFavicon : ...
func GetFavicon(ctx gosip.AuthCnfg) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "favicon")
	}
}
