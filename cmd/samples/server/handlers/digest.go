package handlers

import (
	"fmt"
	"net/http"

	"github.com/koltyakov/gosip"
)

// GetDigest : ...
func GetDigest(ctx gosip.AuthCnfg) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		client := &gosip.SPClient{
			AuthCnfg: ctx,
		}

		digest, err := gosip.GetDigest(client)
		if err != nil {
			fmt.Printf("unable to get digest: %v\n", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf("{\"digest\":\"%s\"}", digest)))
	}
}
