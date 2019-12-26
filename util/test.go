package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/JamesClonk/compose-broker/config"
)

func TestConfig(apiURL string) *config.Config {
	return &config.Config{
		SkipSSL:         true,
		LogLevel:        "debug",
		LogTimestamp:    true,
		Username:        "broker",
		Password:        "pw",
		CatalogFilename: "../catalog.yml",
		API: config.API{
			URL:               apiURL,
			Token:             "deadbeef",
			DefaultDatacenter: "gce:europe-west1",
			DefaultAccountID:  "586eab527c65836dde5533e8",
			Retries:           1,
			RetryInterval:     10 * time.Millisecond,
		},
	}
}

type HttpTestCase struct {
	Method string
	Path   string
	Code   int
	Body   string
	Test   func(string)
}

func TestServer(token string, testCases []HttpTestCase) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(token) > 0 {
			if r.Header.Get("Authorization") != "Bearer "+token {
				w.WriteHeader(401)
				_, _ = w.Write([]byte(`{"errors":"invalid_token"}`))
				return
			}
		}

		for _, test := range testCases {
			if len(test.Method) == 0 {
				test.Method = "GET"
			}
			if strings.HasSuffix(r.RequestURI, test.Path) && r.Method == test.Method {
				if test.Test != nil {
					b, _ := ioutil.ReadAll(r.Body)
					test.Test(string(b))
				}

				w.WriteHeader(test.Code)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, test.Body)
				return
			}
		}
	}))
}

func Body(filename string) string {
	data, _ := ioutil.ReadFile(filename)
	return string(data)
}
