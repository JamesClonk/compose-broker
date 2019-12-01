package broker

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JamesClonk/compose-broker/log"
	"github.com/JamesClonk/compose-broker/util"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestBroker_Health(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	NewRouter(util.TestConfig("")).ServeHTTP(rec, req)
	assert.Equal(t, 200, rec.Code)
	assert.Contains(t, rec.Body.String(), `"status": "ok"`)
}

func TestBroker_BasicAuth(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")

	NewRouter(util.TestConfig("")).ServeHTTP(rec, req)
	assert.Equal(t, 200, rec.Code)
	assert.Contains(t, rec.Body.String(), `"status": "ok"`)
}

func TestBroker_BasicAuth_Unauthorized(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	NewRouter(util.TestConfig("")).ServeHTTP(rec, req)
	assert.Equal(t, 401, rec.Code)
	assert.Contains(t, rec.Body.String(), `"error": "Unauthorized"`)
	assert.Contains(t, rec.Body.String(), `"description": "You are not authorized to access this service broker"`)
}
