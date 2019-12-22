package api

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/JamesClonk/compose-broker/log"
	"github.com/JamesClonk/compose-broker/util"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestAPI_Delete(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{"DELETE", "/api", 202, util.Body("../_fixtures/api_example_valid.json"), nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	c := NewClient(util.TestConfig(apiServer.URL))

	body, err := c.Delete("api")
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, body, `"application": "elastic_search"`)
	assert.Contains(t, body, `"version": "2.4.0"`)
	assert.Equal(t, util.Body("../_fixtures/api_example_valid.json"), body)
}

func TestAPI_Delete_WrongStatusCode(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{"DELETE", "/api", 200, util.Body("../_fixtures/api_example_valid.json"), nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	c := NewClient(util.TestConfig(apiServer.URL))

	_, err := c.Delete("api")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), `unexpected status code: 200, could not parse API response:`)
}

func TestAPI_Get(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{"GET", "/api", 200, util.Body("../_fixtures/api_example_valid.json"), nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	c := NewClient(util.TestConfig(apiServer.URL))

	body, err := c.Get("api")
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, body, `"application": "elastic_search"`)
	assert.Contains(t, body, `"version": "2.4.0"`)
	assert.Equal(t, util.Body("../_fixtures/api_example_valid.json"), body)
}

func TestAPI_Get_WrongStatusCode(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{"GET", "/api", 202, util.Body("../_fixtures/api_example_valid.json"), nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	c := NewClient(util.TestConfig(apiServer.URL))

	_, err := c.Get("api")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), `unexpected status code: 202, could not parse API response:`)
}

func TestAPI_Get_InvalidResponse(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{"GET", "/api", 500, util.Body("../_fixtures/api_example_invalid.json"), nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	c := NewClient(util.TestConfig(apiServer.URL))
	c.Retries = 1
	c.RetryInterval = 10 * time.Millisecond

	_, err := c.Get("api")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), `could not parse API response: {error}`)
}

func TestAPI_Get_ErrorResponse(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{"GET", "/api", 500, util.Body("../_fixtures/api_example_error.json"), nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	c := NewClient(util.TestConfig(apiServer.URL))
	c.Retries = 1
	c.RetryInterval = 10 * time.Millisecond

	_, err := c.Get("api")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), `we've encountered a problem!`)
}

func TestAPI_Get_MultipleErrorsResponse(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{"GET", "/api", 500, util.Body("../_fixtures/api_example_errors.json"), nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	c := NewClient(util.TestConfig(apiServer.URL))
	c.Retries = 1
	c.RetryInterval = 10 * time.Millisecond

	_, err := c.Get("api")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), `api_error: mistake!, big time!`)
	assert.Contains(t, err.Error(), `server_error: fatality!`)
}
