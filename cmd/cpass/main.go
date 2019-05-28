package main

import (
	"flag"
	"fmt"

	"github.com/koltyakov/gosip/cpass"
)

func main() {

	var rawSecret string

	flag.StringVar(&rawSecret, "secret", "", "Raw secret string")
	flag.Parse()

	crypt := cpass.Cpass("")
	secret, _ := crypt.Encode(rawSecret)

	fmt.Println(secret)

}
