package main

import (
	"fmt"

	"github.com/koltyakov/gosip/cpass"
)

func main() {
	c := cpass.NewCpass("")
	myString := "\\\"My string\\\""
	fmt.Printf("Original: %s\n", myString)

	encoded, _ := c.Encrypt([]byte(myString))
	fmt.Printf("Encoded: %x\n", encoded)

	decoded, _ := c.Decrypt(encoded)
	fmt.Printf("Decoded: %s\n", decoded)

	encodedStr, _ := c.Encode(myString)
	fmt.Printf("encodedStr: %s\n", encodedStr)

	decodedStr, _ := c.Decode(encodedStr)
	fmt.Printf("decodedStr: %s\n", decodedStr)

}
