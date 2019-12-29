package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/koltyakov/gosip"
)

// GetWeb : ...
func GetWeb(ctx gosip.AuthCnfg) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		client := &gosip.SPClient{
			AuthCnfg: ctx,
		}

		endpoint := ctx.GetSiteURL() + "/_api/web"
		req, err := http.NewRequest("GET", endpoint, nil)
		if err != nil {
			message := fmt.Sprintf("unable to create a request: %v", err)
			http.Error(w, message, http.StatusBadRequest)
			return
		}

		req.Header.Set("Accept", "application/json;odata=verbose")

		fmt.Printf("requesting endpoint: %s\n", endpoint)
		resp, err := client.Execute(req)
		if err != nil {
			message := fmt.Sprintf("unable to request: %v\n", err)
			http.Error(w, message, http.StatusBadRequest)
			return
		}
		defer resp.Body.Close()

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			message := fmt.Sprintf("unable to read a response: %v\n", err)
			http.Error(w, message, http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}
