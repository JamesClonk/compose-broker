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

func TestAPI_GetScalingByDeploymentID(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{"GET", "/deployments/5821fd28a4b549d06e39886d/scalings", 200, util.Body("../_fixtures/api_get_scaling.json"), nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	c := NewClient(util.TestConfig(apiServer.URL))

	scaling, err := c.GetScalingByDeploymentID("5821fd28a4b549d06e39886d")
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
