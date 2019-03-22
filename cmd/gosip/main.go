package main

import (
	"fmt"

	m "github.com/koltyakov/gosip/test/manual"
)

func main() {
	// m.GetAdfsAuthTest()
	// m.GetWapAuthTest()
	// m.GetWapAdfsAuthTest()
	// m.GetOnlineADFSTest()
	client := m.GetNtlmAuthTest()
	resp, err := m.CheckBasicPost(client)
	if err != nil {
		fmt.Printf("error in CheckBasicPost: %v\n", err)
	}
	fmt.Printf("response from CheckBasicPost: %s\n", resp)
}
