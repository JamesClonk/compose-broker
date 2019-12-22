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

func TestAPI_GetRecipe(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{"GET", "/recipes/5821fd28a4b549d06e39886d", 200, util.Body("../_fixtures/api_get_recipe.json"), nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	c := NewClient(util.TestConfig(apiServer.URL))

	recipe, err := c.GetRecipe("5821fd28a4b549d06e39886d")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "5821fd28a4b549d06e39886d", recipe.ID)
	assert.Equal(t, "Provision", recipe.Name)
	assert.Equal(t, "Recipes::Deployment::Run", recipe.Template)
	assert.Equal(t, "complete", recipe.Status)
	assert.Equal(t, "All operations have completed successfully!", recipe.StatusDetail)
	assert.Equal(t, "5854017d89d50f424e000002", recipe.AccountID)
	assert.Equal(t, "5854017e89d50f424e000192", recipe.DeploymentID)
}

func TestAPI_GetRecipesByDeploymentID(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{"GET", "/deployments/5854017e89d50f424e000192/recipes", 200, util.Body("../_fixtures/api_get_recipes.json"), nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	c := NewClient(util.TestConfig(apiServer.URL))

	recipes, err := c.GetRecipesByDeploymentID("5854017e89d50f424e000192")
	if err != nil {
		t.Fatal(err)
	}

	type deployments struct {
		Embedded struct {
			Recipes Recipes `json:"recipes"`
		} `json:"_embedded"`
	}
	rs := deployments{}
	if err := json.Unmarshal([]byte(util.Body("../_fixtures/api_get_recipes.json")), &rs); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, rs.Embedded.Recipes, recipes)
	assert.Equal(t, "Recipes::Deployment::Run", recipes[0].Template)
	assert.Equal(t, "Recipes::Deployment::Deprovision", recipes[1].Template)
	assert.Equal(t, "Provision", recipes[0].Name)
	assert.Equal(t, "waiting", recipes[1].Status)
	assert.Equal(t, 14, recipes[0].OperationsComplete)
	assert.Equal(t, 0, recipes[1].OperationsTotal)
}
