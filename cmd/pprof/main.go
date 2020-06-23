package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

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

	sp := api.NewSP(client)

	go func() {
		for {
			runner(sp)
			time.Sleep(1 * time.Second)
		}
	}()

	// curl -sK -v http://localhost:6868/debug/pprof/heap > heap.out
	// go tool pprof heap.out
	if err := http.ListenAndServe("localhost:6868", nil); err != nil {
		fmt.Printf("error starting pprof server: %s\n", err)
	}

}

func runner(sp *api.SP) {
	r, err := sp.Web().Select("Title").Get()
	if err != nil {
		fmt.Println(err)
	}
	_ = r.Data().Title
	fmt.Print(".")
}
