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

func TestAPI_GetScaling(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{Method: "GET", Path: "/deployments/5821fd28a4b549d06e39886d/scalings", Code: 200, Body: util.Body("../_fixtures/api_get_scaling.json"), Test: nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	c := NewClient(util.TestConfig(apiServer.URL))

	scaling, err := c.GetScaling("5821fd28a4b549d06e39886d")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 4, scaling.AllocatedUnits)
	assert.Equal(t, 3, scaling.UsedUnits)
	assert.Equal(t, 2, scaling.StartingUnits)
	assert.Equal(t, 1, scaling.MinimumUnits)
	assert.Equal(t, 2048, scaling.MemoryPerUnitInMB)
	assert.Equal(t, 4096, scaling.StoragePerUnitInMB)
	assert.Equal(t, 1024, scaling.UnitSizeInMB)
	assert.Equal(t, "memory", scaling.UnitType)
}

func TestAPI_UpdateScaling(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{Method: "POST", Path: "/deployments/5854017e89d50f424e000192/scalings", Code: 200, Body: util.Body("../_fixtures/api_update_scaling.json"), Test: func(body string) {
			assert.Contains(t, body, `{"deployment":{"units":7}}`)
		}},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	c := NewClient(util.TestConfig(apiServer.URL))

	recipe, err := c.UpdateScaling("5854017e89d50f424e000192", 7)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "570bcb3fee4cde000e000002", recipe.ID)
	assert.Equal(t, "Scale deployment to 7 units", recipe.Name)
	assert.Equal(t, "Recipes::Deployment::Run", recipe.Template)
	assert.Equal(t, "complete", recipe.Status)
	assert.Equal(t, "All operations have completed successfully!", recipe.StatusDetail)
	assert.Equal(t, "586eab527c65836dde5533e8", recipe.AccountID)
	assert.Equal(t, "5854017e89d50f424e000192", recipe.DeploymentID)
}
