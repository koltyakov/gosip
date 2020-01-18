package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/koltyakov/gosip/api"
	m "github.com/koltyakov/gosip/test/manual"
)

func main() {

	strategy := flag.String("strategy", "saml", "Auth strategy code")
	flag.Parse()

	client, err := m.GetTestClient(*strategy)
	if err != nil {
		log.Fatal(err)
	}

	// Manual test code is below

	sp := api.NewSP(client)
	res, err := sp.Web().Select("Title").Get()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", res.Data().Title)

}
