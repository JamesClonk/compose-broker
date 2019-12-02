package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/JamesClonk/compose-broker/config"
	"github.com/JamesClonk/compose-broker/log"
	"github.com/JamesClonk/compose-broker/util"
	"github.com/parnurzeal/gorequest"
)

type Client struct {
	Config           config.API
	Mutex            *sync.Mutex
	HTTPClient       *http.Client
	Retries          int
	RetryInterval    time.Duration
	RetryStatusCodes []int
}

func NewClient(c *config.Config) *Client {
	client := &Client{
		Config:        c.API,
		Mutex:         &sync.Mutex{},
		HTTPClient:    util.NewHttpClient(c),
		Retries:       3,
		RetryInterval: 3 * time.Second,
		RetryStatusCodes: []int{
			http.StatusRequestTimeout,
			http.StatusTooManyRequests,
			http.StatusBadGateway,
			http.StatusServiceUnavailable,
			http.StatusGatewayTimeout,
		},
	}
	return client
}

func (c *Client) newRequest(method, endpoint string) *gorequest.SuperAgent {
	targetURL := fmt.Sprintf("%s/%s", c.Config.URL, endpoint)
	log.Debugf("Compose.io API HTTP request [%v:%v]", method, targetURL)

	return gorequest.New().
		CustomMethod(method, targetURL).
		Set("Authorization", "Bearer "+c.Config.Token).
		Set("Content-type", "application/json; charset=utf-8").
		Retry(c.Retries, c.RetryInterval, c.RetryStatusCodes...)
}

func (c *Client) GetJSON(endpoint string) (string, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	response, body, errs := c.newRequest("GET", endpoint).End()
	if response.StatusCode != 200 {
		errs = composeErrors(response.StatusCode, body)
	}

	var err error
	if len(errs) > 0 {
		var errorMessage string
		for _, err := range errs {
			if len(errorMessage) > 0 {
				errorMessage = fmt.Sprintf("%s: %s", errorMessage, err.Error())
			} else {
				errorMessage = err.Error()
			}
		}
		err = errors.New(errorMessage)
	}
	return body, err
}

func composeErrors(statuscode int, body string) []error {
	// Compose.io error types
	type errors struct {
		Error map[string][]string `json:"errors,omitempty"`
	}
	type simpleError struct {
		Error string `json:"errors"`
	}

	errs := []error{}
	myerrors := errors{}

	err := json.Unmarshal([]byte(body), &myerrors)
	if err != nil {
		simpleerror := simpleError{}
		err := json.Unmarshal([]byte(body), &simpleerror)
		if err != nil {
			errs = append(errs, fmt.Errorf("status code %d: %s", statuscode, body))
		} else {
			errs = append(errs, fmt.Errorf("%s", simpleerror.Error))
		}
	} else {
		errs = append(errs, fmt.Errorf("%v", myerrors.Error))
	}
	return errs
}
