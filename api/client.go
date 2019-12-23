package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
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
		Retries:       c.API.Retries,
		RetryInterval: c.API.RetryInterval,
		RetryStatusCodes: []int{
			http.StatusRequestTimeout,
			http.StatusTooManyRequests,
			http.StatusBadGateway,
			http.StatusServiceUnavailable,
			http.StatusGatewayTimeout,
			http.StatusInternalServerError,
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

func (c *Client) Get(endpoint string) (string, error) {
	return c.Do("GET", endpoint, "", 200)
}

func (c *Client) Post(endpoint, payload string) (string, error) {
	return c.Do("POST", endpoint, payload, 200)
}

func (c *Client) PostAsync(endpoint, payload string) (string, error) {
	return c.Do("POST", endpoint, payload, 202)
}

func (c *Client) Delete(endpoint string) (string, error) {
	return c.Do("DELETE", endpoint, "", 202)
}

func (c *Client) Do(method, endpoint, payload string, code int) (string, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	request := c.newRequest(method, endpoint)
	if method == "POST" {
		request = request.Send(payload)
	}
	response, body, errs := request.End()
	if response.StatusCode != code {
		errs = append(errs, fmt.Errorf("unexpected status code: %d", response.StatusCode))
		errs = append(errs, composeErrors(body)...)
	}

	var err error
	if len(errs) > 0 {
		var errorMessage string
		for _, err := range errs {
			if len(errorMessage) > 0 {
				errorMessage = fmt.Sprintf("%s, %s", errorMessage, err.Error())
			} else {
				errorMessage = err.Error()
			}
		}
		err = errors.New(errorMessage)
	}
	return body, err
}

func composeErrors(body string) []error {
	errs := make([]error, 0)

	// Compose.io error types
	type multiErrors struct {
		Error map[string][]string `json:"errors,omitempty"`
	}
	type singleError struct {
		Error string `json:"errors"`
	}

	multi := multiErrors{}
	if err := json.Unmarshal([]byte(body), &multi); err != nil {
		single := singleError{}
		if err := json.Unmarshal([]byte(body), &single); err != nil {
			errs = append(errs, fmt.Errorf("could not parse API response: %s", body))
		} else {
			errs = append(errs, fmt.Errorf("%s", single.Error))
		}
	} else {
		for key, value := range multi.Error {
			errs = append(errs, fmt.Errorf("%s: %s", key, strings.Join(value, ", ")))
		}
	}
	return errs
}
