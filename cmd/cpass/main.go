package main

import (
	"flag"
	"fmt"

	"github.com/koltyakov/gosip/cpass"
)

func main() {

	var rawSecret string
	var masterKey string
	var mode string

	flag.StringVar(&rawSecret, "secret", "", "Raw secret string")
	flag.StringVar(&masterKey, "master", "", "Master key string")
	flag.StringVar(&mode, "mode", "encode", "Mode: encode/decode")
	flag.Parse()

	crypt := cpass.Cpass("")

	if mode == "encode" {
		secret, _ := crypt.Encode(rawSecret)
		fmt.Println(secret)
	}

	if mode == "decode" {
		secret, _ := crypt.Decode(rawSecret)
		fmt.Println(secret)
	}

}
