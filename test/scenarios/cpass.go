package scenarios

import (
	"fmt"

	"github.com/koltyakov/gosip/cpass"
)

// CpassDummyTest : test scenario
func CpassDummyTest() {
	c := cpass.Cpass("DUMMY_KEY")
	myString := "secret"
	dummySecret, _ := c.Encode(myString)
	fmt.Printf("dummySecret: %s\n", dummySecret)
	decodedStr, _ := c.Decode(dummySecret)
	fmt.Printf("decodedStr: %s\n", decodedStr)
}

// CpassAutoModeTest : test scenario
func CpassAutoModeTest() {
	c := cpass.Cpass("")
	myString := "secret"
	fmt.Printf("Original: %s\n", myString)
	encodedStr, _ := c.Encode(myString)
	fmt.Printf("encodedStr: %s\n", encodedStr)
	decodedStr, _ := c.Decode(encodedStr)
	fmt.Printf("decodedStr: %s\n", decodedStr)

	// cpassFromNodeJS, _ := c.Decode("eefd0b898b9c9aa80b6ced46204ba228JpCTeYtSJSBjdxZrRk24kg==")
	// fmt.Printf("cpassFromNodeJS: %s\n", cpassFromNodeJS)
}
