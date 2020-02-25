package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/koltyakov/gosip/api"
	m "github.com/koltyakov/gosip/test/manual"
)

func main() {

	strategy := flag.String("strategy", "fba", "Auth strategy code")
	flag.Parse()

	client, err := m.GetTestClient(*strategy)
	if err != nil {
		log.Fatal(err)
	}

	// Define requests hook handlers
	// client.Hooks = &gosip.HookHandlers{
	// 	OnError: func(e *gosip.HookEvent) {
	// 		fmt.Println("\n======= On Error ========")
	// 		fmt.Printf("URL: %s\n", e.Request.URL)
	// 		fmt.Printf("StatusCode: %d\n", e.StatusCode)
	// 		fmt.Printf("Error: %s\n", e.Error)
	// 		fmt.Printf("took %f seconds\n", time.Since(e.StartedAt).Seconds())
	// 		fmt.Printf("=========================\n\n")
	// 	},
	// 	OnRetry: func(e *gosip.HookEvent) {
	// 		fmt.Println("\n======= On Retry ========")
	// 		fmt.Printf("URL: %s\n", e.Request.URL)
	// 		fmt.Printf("StatusCode: %d\n", e.StatusCode)
	// 		fmt.Printf("Error: %s\n", e.Error)
	// 		fmt.Printf("took %f seconds\n", time.Since(e.StartedAt).Seconds())
	// 		fmt.Printf("=========================\n\n")
	// 	},
	// 	OnRequest: func(e *gosip.HookEvent) {
	// 		if e.Error == nil {
	// 			fmt.Println("\n====== On Request =======")
	// 			fmt.Printf("URL: %s\n", e.Request.URL)
	// 			fmt.Printf("auth injection took %f seconds\n", time.Since(e.StartedAt).Seconds())
	// 			fmt.Printf("=========================\n\n")
	// 		}
	// 	},
	// 	OnResponse: func(e *gosip.HookEvent) {
	// 		if e.Error == nil {
	// 			fmt.Println("\n====== On Response =======")
	// 			fmt.Printf("URL: %s\n", e.Request.URL)
	// 			fmt.Printf("StatusCode: %d\n", e.StatusCode)
	// 			fmt.Printf("took %f seconds\n", time.Since(e.StartedAt).Seconds())
	// 			fmt.Printf("==========================\n\n")
	// 		}
	// 	},
	// }

	// Manual test code is below

	sp := api.NewSP(client)
	res, err := sp.Web().Select("Title").Get()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", res.Data().Title)

	// if _, err := sp.Web().Lists().GetByTitle("NotExisting").Get(); err != nil {
	// 	log.Fatal(err)
	// }

}
