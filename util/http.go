package util

import (
	"crypto/tls"
	"net/http"
	"os"
	"time"

	"github.com/JamesClonk/compose-broker/config"
)

func NewHttpClient(c *config.Config) *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: c.SkipSSL,
		},
	}

	proxy := os.Getenv("http_proxy")
	if len(proxy) > 0 {
		tr.Proxy = http.ProxyFromEnvironment
	}

	return &http.Client{
		Timeout:   time.Duration(33 * time.Second),
		Transport: tr,
	}
}
