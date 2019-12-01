package api

import (
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/JamesClonk/compose-broker/config"
	"github.com/JamesClonk/compose-broker/log"
	"github.com/JamesClonk/compose-broker/util"
)

type Client struct {
	API    config.API
	Mutex  *sync.Mutex
	Client *http.Client
}

func NewClient(c *config.Config) *Client {
	client := &Client{
		API:    c.API,
		Mutex:  &sync.Mutex{},
		Client: util.NewHttpClient(c),
	}
	return client
}

// TODO: switch to https://github.com/parnurzeal/gorequest
// see https://github.com/compose/gocomposeapi/blob/master/composeapi.go example..
func (c *Client) Do(req *http.Request) (int, []byte, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	log.Debugf("Compose.io API request [%v:%v]", req.Method, req.URL.RequestURI())

	if len(c.API.Token) > 0 {
		// TODO:-H "Authorization: Bearer [[app:Authorization]]"
		req.Header.Set("Authorization", "Bearer "+c.API.Token)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := c.Client.Do(req)
	if err != nil {
		return 500, nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, err
	}
	return resp.StatusCode, body, nil
}
