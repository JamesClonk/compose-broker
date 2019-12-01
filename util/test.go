package util

import (
	"crypto/subtle"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

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
			DefaultWhitelist:  "0.0.0.0/0",
		},
	}
}

type HttpTestCase struct {
	HttpStatusCode int
	ResponseBody   string
	TestFunc       func(string)
}

func TestServer(username, password string, testCases map[string]HttpTestCase) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(username) > 0 || len(password) > 0 {
			user, pw, ok := r.BasicAuth()
			if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pw), []byte(password)) != 1 {
				w.Header().Set("WWW-Authenticate", `Basic realm="test"`)
				w.WriteHeader(401)
				w.Write([]byte("Not authorized"))
				return
			}
		}

		for path, test := range testCases {
			var method string
			if strings.Contains(path, "::") {
				values := strings.SplitN(path, "::", 2)
				path = values[1]
				method = values[0]
			}
			if strings.HasSuffix(r.RequestURI, path) &&
				(len(method) == 0 || method == r.Method) {
				if test.TestFunc != nil {
					b, _ := ioutil.ReadAll(r.Body)
					test.TestFunc(string(b))
				}

				w.WriteHeader(test.HttpStatusCode)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, test.ResponseBody)
				return
			}
		}
	}))
}

func Body(filename string) string {
	data, _ := ioutil.ReadFile(filename)
	return string(data)
}
