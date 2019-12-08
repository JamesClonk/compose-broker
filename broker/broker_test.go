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

func TestBroker_Write(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/yolo", nil)
	if err != nil {
		t.Fatal(err)
	}

	b := NewBroker(util.TestConfig(""))
	b.write(rec, req, 418, map[string]string{"text": "example"})

	assert.Equal(t, 418, rec.Code)
	assert.Equal(t, "compose-broker", rec.Header().Get("X-Compose-Broker"))
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Equal(t, `{
  "text": "example"
}`, rec.Body.String())
}

func TestBroker_Error(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/wrong", nil)
	if err != nil {
		t.Fatal(err)
	}

	b := NewBroker(util.TestConfig(""))
	b.Error(rec, req, 406, "Wrong", "You are wrong!")

	assert.Equal(t, 406, rec.Code)
	assert.Equal(t, "compose-broker", rec.Header().Get("X-Compose-Broker"))
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Equal(t, `{
  "description": "You are wrong!",
  "error": "Wrong"
}`, rec.Body.String())
}
