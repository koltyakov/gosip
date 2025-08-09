package ondemand

import "github.com/koltyakov/lorca"

// dimension struct for defining height and width
type dimension struct {
	Width  int
	Height int
}

// screenSize cache
var screenSize = dimension{}

// getScreenSize gets screen width and height
func getScreenSize() (dimension, error) {
	if screenSize.Width != 0 && screenSize.Height != 0 {
		return screenSize, nil
	}

	// Open Chrome in headless mode
	ui, err := lorca.New("data:text/html,", "", 0, 0, "--remote-allow-origins=*") //, "--headless")
	if err != nil {
		// Can't open Chrome, likely is not installed
		return screenSize, err
	}

	defer func() { _ = ui.Close() }()

	screenSize.Width = ui.Eval("screen.availWidth").Int()
	screenSize.Height = ui.Eval("screen.availHeight").Int()

	return screenSize, nil
}
