package api

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/JamesClonk/compose-broker/log"
	"github.com/JamesClonk/compose-broker/util"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestAPI_GetDatabases(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{"GET", "/databases", 200, util.Body("../_fixtures/api_get_databases.json"), nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	c := NewClient(util.TestConfig(apiServer.URL))

	dbs, err := c.GetDatabases()
	if err != nil {
		t.Fatal(err)
	}

	ds := struct {
		Embedded struct {
			Databases Databases `json:"applications"`
		} `json:"_embedded"`
	}{}
	if err := json.Unmarshal([]byte(util.Body("../_fixtures/api_get_databases.json")), &ds); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, ds.Embedded.Databases, dbs)
}
