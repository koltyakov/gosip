package main

import (
	m "github.com/koltyakov/gosip/test/manual"
)

func main() {
	m.GetAdfsAuthTest()
	m.GetWapAuthTest()
	m.GetWapAdfsAuthTest()
}
