package api

import (
	"io/ioutil"
	"testing"

	"github.com/JamesClonk/compose-broker/log"
	"github.com/JamesClonk/compose-broker/util"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestAPI_GetAccounts(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{"GET", "/accounts", 200, util.Body("../_fixtures/api_get_accounts.json"), nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	c := NewClient(util.TestConfig(apiServer.URL))

	accs, err := c.GetAccounts()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "52a50cb96230650018000000", accs[0].ID)
	assert.Equal(t, "Northwind", accs[0].Name)
	assert.Equal(t, "northwind", accs[0].Slug)
}
