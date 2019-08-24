package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/koltyakov/gosip/auth/adfs"
)

func main() {

	var siteURL string
	var username string
	var password string
	var relyingParty string
	var adfsURL string
	var adfsCookie string
	var configPath string
	var outFormat string

	flag.StringVar(&siteURL, "siteUrl", "", "SharePoint Site Url")
	flag.StringVar(&username, "username", "", "User login")
	flag.StringVar(&password, "password", "", "User password")
	flag.StringVar(&relyingParty, "relyingParty", "", "Relying party")
	flag.StringVar(&adfsURL, "adfsUrl", "", "ADFS Url")
	flag.StringVar(&adfsCookie, "adfsCookie", "", "ADFS Cookie")
	flag.StringVar(&configPath, "configPath", "", "Connection config path")
	flag.StringVar(&outFormat, "outFormat", "json", "Output Format: raw | json")

	flag.Parse()

	auth := &adfs.AuthCnfg{
		SiteURL:      siteURL,
		Username:     username,
		Password:     password,
		RelyingParty: relyingParty,
		AdfsURL:      adfsURL,
		AdfsCookie:   adfsCookie,
	}

	if configPath != "" {
		err := auth.ReadConfig(configPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to read config: %v", err)
			os.Exit(1)
		}
	}

	authCookie, err := auth.GetAuth()
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to authenticate: %v", err)
		os.Exit(1)
	}

	if authCookie == "" {
		fmt.Fprint(os.Stderr, "can't get auth cookie")
		os.Exit(1)
	}

	if outFormat == "raw" {
		fmt.Fprint(os.Stdout, authCookie)
		os.Exit(0)
	}

	cookies := strings.Split(authCookie, "; ")

	json := "{"
	for i, cookie := range cookies {
		c := strings.SplitN(cookie, "=", 2)
		json += fmt.Sprintf("\"%s\":\"%s\"", c[0], c[1])
		if i+1 < len(cookies) {
			json += ","
		}
	}
	json += "}"

	fmt.Fprint(os.Stdout, json)
}
