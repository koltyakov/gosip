package api

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/koltyakov/gosip"
	a "github.com/koltyakov/gosip/auth"
	"github.com/koltyakov/gosip/auth/ntlm"
	"github.com/koltyakov/gosip/auth/saml"
	h "github.com/koltyakov/gosip/test/helpers"
)

var (
	ci         bool
	skip       bool
	heavyTests bool
	envCode    string
	spClient   *gosip.SPClient
	headers    struct {
		verbose         *RequestConfig
		minimalmetadata *RequestConfig
		nometadata      *RequestConfig
	}
)

// Request counters
var requestCntrs = struct {
	Errors    int32
	Responses int32
	Retries   int32
	Requests  int32
}{
	Errors:    0,
	Responses: 0,
	Retries:   0,
	Requests:  0,
}

func init() {
	readDotEnv()

	ci = os.Getenv("SPAUTH_CI") == "true"
	skip = os.Getenv("SPAPI_SKIP_TESTS") == "true"
	heavyTests = os.Getenv("SPAPI_HEAVY_TESTS") == "true"
	envCode = os.Getenv("SPAUTH_ENVCODE")

	if envCode == "" && !ci && !skip {
		envCode = "spo"
	}

	envResolver := map[string]func() *gosip.SPClient{
		"spo": func() *gosip.SPClient {
			cnfgPath := "./config/integration/private.spo.json"
			if ci {
				auth := &saml.AuthCnfg{}
				auth.SiteURL = os.Getenv("ENV_SPO_SITEURL")
				auth.Username = os.Getenv("ENV_SPO_USERNAME")
				auth.Password = os.Getenv("ENV_SPO_PASSWORD")
				if err := h.CheckAuthProps(auth, []string{"SiteURL", "Username", "Password"}); err != nil {
					return nil
				}
				client := &gosip.SPClient{AuthCnfg: auth}
				// Pre-auth for tests not include auth timing involved
				if _, _, err := client.AuthCnfg.GetAuth(); err != nil {
					fmt.Printf("can't auth, %s\n", err)
					// Force all test being skipped in case of auth issues
					return nil
				}
				return client
			}

			authCnfg, err := a.NewAuthFromFile(resolveCnfgPath(cnfgPath))
			if err != nil {
				return nil
			}
			client := &gosip.SPClient{AuthCnfg: authCnfg}
			// Pre-auth for tests not include auth timing involved
			if _, _, err := client.AuthCnfg.GetAuth(); err != nil {
				fmt.Printf("can't auth, %s\n", err)
				// Force all test being skipped in case of auth issues
				return nil
			}
			return client
		},
		"2013": func() *gosip.SPClient {
			cnfgPath := "./config/integration/private.2013.json"
			auth := &ntlm.AuthCnfg{}
			if err := auth.ReadConfig(resolveCnfgPath(cnfgPath)); err != nil {
				return nil
			}
			if err := h.CheckAuthProps(auth, []string{"SiteURL", "Username", "Password", "Domain"}); err != nil {
				return nil
			}
			client := &gosip.SPClient{AuthCnfg: auth}
			// Pre-auth for tests not include auth timing involved
			if _, _, err := client.AuthCnfg.GetAuth(); err != nil {
				fmt.Printf("can't auth, %s\n", err)
				// Force all test being skipped in case of auth issues
				return nil
			}
			return client
		},
	}

	handlers := &gosip.HookHandlers{
		OnError: func(e *gosip.HookEvent) {
			atomic.AddInt32(&requestCntrs.Errors, 1)
		},
		OnResponse: func(e *gosip.HookEvent) {
			atomic.AddInt32(&requestCntrs.Responses, 1)
		},
		OnRetry: func(e *gosip.HookEvent) {
			atomic.AddInt32(&requestCntrs.Retries, 1)
		},
		OnRequest: func(e *gosip.HookEvent) {
			atomic.AddInt32(&requestCntrs.Requests, 1)
		},
	}

	if envCode != "" && envResolver[envCode] != nil {
		spClient = envResolver[envCode]()
		if spClient == nil {
			skip = true
			fmt.Printf("Warning: can't resolve auth context for %s\n", envCode)
			return
		}

		spClient.Hooks = handlers

		spClient.RetryPolicies = map[int]int{
			404: 1,
		}
	}

	spClient.Timeout = 30 * time.Second

	setHeadersPresets()
}

func resolveCnfgPath(relativePath string) string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Join(path.Dir(filename), "..", relativePath)
}

func checkClient(t *testing.T) {
	if spClient == nil {
		t.Skip("no auth context provided")
	}
	if os.Getenv("GOMAXPROCS") != "" {
		t.Parallel()
	}
}

func setHeadersPresets() {
	headers.verbose = HeadersPresets.Verbose
	headers.minimalmetadata = HeadersPresets.Minimalmetadata
	headers.nometadata = HeadersPresets.Nometadata
}

func readDotEnv() {
	envFilePath := resolveCnfgPath(".env")
	envFile, err := os.Open(envFilePath)
	if err != nil {
		return
	}
	defer func() { _ = envFile.Close() }()

	byteValue, _ := io.ReadAll(envFile)
	keyVals := strings.Split(string(byteValue), "\n")
	for _, keyVal := range keyVals {
		kv := strings.SplitN(keyVal, "=", 2)
		if len(kv) == 2 {
			_ = os.Setenv(kv[0], kv[1])
		}
	}
}
