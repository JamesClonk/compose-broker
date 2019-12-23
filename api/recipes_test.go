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

func TestAPI_GetRecipes(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{"GET", "/deployments/5854017e89d50f424e000192/recipes", 200, util.Body("../_fixtures/api_get_recipes.json"), nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	c := NewClient(util.TestConfig(apiServer.URL))

	recipes, err := c.GetRecipes("5854017e89d50f424e000192")
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
	assert.NotEqual(t, rs.Embedded.Recipes, recipes) // returned by being sorted by "updated_at desc"

	assert.Equal(t, "Recipes::Deployment::Run", recipes[1].Template)
	assert.Equal(t, "Recipes::Deployment::Deprovision", recipes[0].Template)
	assert.Equal(t, "Provision", recipes[1].Name)
	assert.Equal(t, "Deprovision", recipes[0].Name)
	assert.Equal(t, "complete", recipes[1].Status)
	assert.Equal(t, "waiting", recipes[0].Status)
	assert.Equal(t, 14, recipes[1].OperationsComplete)
	assert.Equal(t, 0, recipes[0].OperationsTotal)

	recipes.SortByCreatedAt()
	assert.Equal(t, rs.Embedded.Recipes, recipes)
	assert.Equal(t, "Provision", recipes[0].Name)
	assert.Equal(t, "Deprovision", recipes[1].Name)

	recipes.SortByUpdatedAt()
	assert.Equal(t, "Provision", recipes[1].Name)
	assert.Equal(t, "Deprovision", recipes[0].Name)
}
