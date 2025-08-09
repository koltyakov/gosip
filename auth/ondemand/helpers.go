package ondemand

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/koltyakov/lorca"
)

// Settings dialog window size
var dlg = &dimension{
	Width:  580,
	Height: 530,
}

// onDemandAuthFlow authenticates using On-Demand flow
func (c *AuthCnfg) onDemandAuthFlow(initialCookies *Cookies) (*Cookies, error) {
	var args []string

	if screen, err := getScreenSize(); err == nil && screen.Height != 0 && screen.Width != 0 {
		args = append(
			args,
			fmt.Sprintf(
				"--window-position=%d,%d",
				(screen.Width-dlg.Width)/2,
				(screen.Height-dlg.Height)/2,
			),
			"--remote-allow-origins=*", // #6
		)
	}

	if c.ChromeArgs != nil {
		for arg, val := range *c.ChromeArgs {
			// Arg must start with "--" to be a valid Chrome argument
			if !strings.HasPrefix(arg, "--") {
				args = append(args, fmt.Sprintf("--%s=%s", arg, val))
				fmt.Printf("Adding Chrome arg: %s=%s\n", arg, val)
			}
		}
	}

	startURL := fmt.Sprintf("data:text/html;base64,%s", base64.StdEncoding.EncodeToString([]byte(getStartHTML(c.SiteURL))))
	ui, err := lorca.New(startURL, "", dlg.Width, dlg.Height, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = ui.Close() }()

	if initialCookies != nil {
		for _, cookie := range *initialCookies {
			ui.Send("Network.setCookie", cookie.toMap())
		}
	}

	_ = ui.Load(c.SiteURL)

	cookies := &Cookies{}
	var e error

	go func() {
		currentURL := ""
		for strings.Index(strings.ToLower(currentURL), strings.ToLower(c.SiteURL)) == -1 {
			newURL := ui.Eval("window.location.href").String()
			if currentURL != newURL {
				currentURL = newURL
			}
			time.Sleep(500 * time.Microsecond)
		}
		resp := ui.Send("Network.getCookies", nil)
		if resp.Err() != nil {
			e = resp.Err()
			return
		}
		if err := resp.Object()["cookies"].To(&cookies); err != nil {
			e = err
		}
		_ = ui.Close()
	}()

	<-ui.Done()

	if len(*cookies) == 0 {
		e = fmt.Errorf("can't get authentication cookies")
	}

	return cookies, e
}

func getStartHTML(siteURL string) string {
	return `
		<html>
			<head lang="en">
				<title>Connecting to site: ` + siteURL + `</title>
				<style>
					body {
						background-color: #4b4242;
						font-family: "Segoe UI", "Segoe UI Web (West European)", "Segoe UI", -apple-system, BlinkMacSystemFont, Roboto, "Helvetica Neue", sans-serif;
						text-align: center;
					}
					.header {
						color: #fff;
						margin-top: 160px;
					}
					svg {
						height: 90px;
					}
				</style>
			</head>
			<body>
				<h1 class="header">Connecting to site</h1>
				<svg version="1.1" id="L7" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" x="0px" y="0px" viewBox="0 0 100 100" enable-background="new 0 0 100 100" xml:space="preserve">
					<path fill="#fff" d="M31.6,3.5C5.9,13.6-6.6,42.7,3.5,68.4c10.1,25.7,39.2,38.3,64.9,28.1l-3.1-7.9c-21.3,8.4-45.4-2-53.8-23.3c-8.4-21.3,2-45.4,23.3-53.8L31.6,3.5z">
						<animateTransform attributeName="transform" attributeType="XML" type="rotate"dur="2s" from="0 50 50"to="360 50 50" repeatCount="indefinite" />
					</path>
					<path fill="#fff" d="M42.3,39.6c5.7-4.3,13.9-3.1,18.1,2.7c4.3,5.7,3.1,13.9-2.7,18.1l4.1,5.5c8.8-6.5,10.6-19,4.1-27.7c-6.5-8.8-19-10.6-27.7-4.1L42.3,39.6z">
						<animateTransform attributeName="transform" attributeType="XML" type="rotate"dur="1s" from="0 50 50"to="-360 50 50" repeatCount="indefinite" />
					</path>
					<path fill="#fff" d="M82,35.7C74.1,18,53.4,10.1,35.7,18S10.1,46.6,18,64.3l7.6-3.4c-6-13.5,0-29.3,13.5-35.3s29.3,0,35.3,13.5L82,35.7z">
						<animateTransform attributeName="transform" attributeType="XML" type="rotate"dur="2s" from="0 50 50"to="360 50 50" repeatCount="indefinite" />
					</path>
				</svg>
			</body>
		</html>
	`
}
