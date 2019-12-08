package api

import (
	"io/ioutil"
	"log"
	"testing"
	"time"

	"github.com/JamesClonk/compose-broker/util"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestAPI_GetJSON(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{"GET", "/api", 200, util.Body("../_fixtures/api_example_valid.json"), nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	c := NewClient(util.TestConfig(apiServer.URL))

	body, err := c.GetJSON("api")
	if err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, body, `"application": "elastic_search"`)
	assert.Contains(t, body, `"version": "2.4.0"`)
	assert.Equal(t, util.Body("../_fixtures/api_example_valid.json"), body)
}

func TestAPI_GetJSON_Invalid(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{"GET", "/api", 500, util.Body("../_fixtures/api_example_invalid.json"), nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	c := NewClient(util.TestConfig(apiServer.URL))
	c.Retries = 1
	c.RetryInterval = 10 * time.Millisecond

	_, err := c.GetJSON("api")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), `could not parse API response: {error}`)
}

func TestAPI_GetJSON_Error(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{"GET", "/api", 500, util.Body("../_fixtures/api_example_error.json"), nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	c := NewClient(util.TestConfig(apiServer.URL))
	c.Retries = 1
	c.RetryInterval = 10 * time.Millisecond

	_, err := c.GetJSON("api")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), `we've encountered a problem!`)
}

func TestAPI_GetJSON_Errors(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{"GET", "/api", 500, util.Body("../_fixtures/api_example_errors.json"), nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	c := NewClient(util.TestConfig(apiServer.URL))
	c.Retries = 1
	c.RetryInterval = 10 * time.Millisecond

	_, err := c.GetJSON("api")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), `api_error: mistake!, big time!`)
	assert.Contains(t, err.Error(), `server_error: fatality!`)
}
