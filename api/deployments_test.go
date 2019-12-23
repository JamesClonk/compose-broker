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

func TestAPI_GetDeployments(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{"GET", "/deployments", 200, util.Body("../_fixtures/api_get_deployments.json"), nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	c := NewClient(util.TestConfig(apiServer.URL))

	dbs, err := c.GetDeployments()
	if err != nil {
		t.Fatal(err)
	}

	type deployments struct {
		Embedded struct {
			Deployments Deployments `json:"deployments"`
		} `json:"_embedded"`
	}
	ds := deployments{}
	if err := json.Unmarshal([]byte(util.Body("../_fixtures/api_get_deployments.json")), &ds); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, ds.Embedded.Deployments, dbs)
	assert.Equal(t, "fizz-production", dbs[0].Name)
	assert.Equal(t, "bill-to-test-team", dbs[0].CustomerBillingCode)
	assert.Equal(t, "https://app.compose.io/northwind/deployments/test-deployment-2{?embed}", dbs[1].Links.ComposeWebUI.HREF)
	assert.Equal(t, "redis", dbs[1].Type)
}

func TestAPI_GetDeployment(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{"GET", "/deployments/5854017e89d50f424e000192", 200, util.Body("../_fixtures/api_get_deployment.json"), nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	c := NewClient(util.TestConfig(apiServer.URL))

	deployment, err := c.GetDeployment("5854017e89d50f424e000192")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "5854017e89d50f424e000192", deployment.ID)
	assert.Equal(t, "5854017d89d50f424e000002", deployment.AccountID)
	assert.Equal(t, "fizz-production", deployment.Name)
	assert.Equal(t, "postgresql", deployment.Type)
	assert.Equal(t, "the production fizz db", deployment.Notes)
	assert.Equal(t, "bill-to-fizz", deployment.CustomerBillingCode)
	assert.Equal(t, "59a6a6238a681830479c80f8", deployment.ClusterID)
	assert.Equal(t, "9.6.3", deployment.Version)
	assert.Equal(t, "", deployment.CACertificateBase64)
	assert.Equal(t, "postgres://compose:XXXX@customer-cluster.1.compose.direct:10020/compose", deployment.ConnectionStrings.Direct[0])
	assert.Equal(t, "psql \"sslmode=require host=cpu.blazzleblazzle.compose.direct port=10000 dbname=compose user=compose\"", deployment.ConnectionStrings.CLI[0])
	assert.Equal(t, "https://app.compose.io/northwind/deployments/fizz-production{?embed}", deployment.Links.ComposeWebUI.HREF)
}

func TestAPI_GetDeploymentByName(t *testing.T) {
	getDeploymentByIDCalled := false
	test := []util.HttpTestCase{
		util.HttpTestCase{"GET", "/deployments", 200, util.Body("../_fixtures/api_get_deployments.json"), nil},
		util.HttpTestCase{"GET", "/deployments/5854017e89d50f424e000192", 200, util.Body("../_fixtures/api_get_deployment.json"), func(body string) {
			getDeploymentByIDCalled = true
		}},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	c := NewClient(util.TestConfig(apiServer.URL))

	deployment, err := c.GetDeploymentByName("fizz-production")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, getDeploymentByIDCalled) // should be called to get deployment details
	assert.Equal(t, "5854017e89d50f424e000192", deployment.ID)
	assert.Equal(t, "5854017d89d50f424e000002", deployment.AccountID)
	assert.Equal(t, "fizz-production", deployment.Name)
	assert.Equal(t, "postgresql", deployment.Type)
	assert.Equal(t, "the production fizz db", deployment.Notes)
	assert.Equal(t, "bill-to-fizz", deployment.CustomerBillingCode)
	assert.Equal(t, "59a6a6238a681830479c80f8", deployment.ClusterID)
	assert.Equal(t, "9.6.3", deployment.Version)
	assert.Equal(t, "", deployment.CACertificateBase64)
	assert.Equal(t, "postgres://compose:XXXX@customer-cluster.1.compose.direct:10020/compose", deployment.ConnectionStrings.Direct[0])
	assert.Equal(t, "psql \"sslmode=require host=cpu.blazzleblazzle.compose.direct port=10000 dbname=compose user=compose\"", deployment.ConnectionStrings.CLI[0])
	assert.Equal(t, "https://app.compose.io/northwind/deployments/fizz-production{?embed}", deployment.Links.ComposeWebUI.HREF)
}

func TestAPI_GetDeploymentByName_Unknown(t *testing.T) {
	getDeploymentByIDCalled := false
	test := []util.HttpTestCase{
		util.HttpTestCase{"GET", "/deployments", 200, util.Body("../_fixtures/api_get_deployments.json"), nil},
		util.HttpTestCase{"GET", "/deployments/5854017e89d50f424e000192", 200, util.Body("../_fixtures/api_get_deployment.json"), func(body string) {
			getDeploymentByIDCalled = true
		}},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	c := NewClient(util.TestConfig(apiServer.URL))

	deployment, err := c.GetDeploymentByName("invalid name")
	assert.Error(t, err)
	assert.Empty(t, deployment)
	assert.Equal(t, false, getDeploymentByIDCalled) // should not be called, since deployment name could not be found
}

func TestAPI_DeleteDeployment(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{"GET", "/deployments/5854017e89d50f424e000192", 200, util.Body("../_fixtures/api_get_deployment.json"), nil},
		util.HttpTestCase{"DELETE", "/deployments/5854017e89d50f424e000192", 202, util.Body("../_fixtures/api_delete_deployment.json"), nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	c := NewClient(util.TestConfig(apiServer.URL))

	recipe, err := c.DeleteDeployment("5854017e89d50f424e000192")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "570bf60a70ea13000d000000", recipe.ID)
	assert.Equal(t, "Deprovision", recipe.Name)
	assert.Equal(t, "Recipes::Deployment::Deprovision", recipe.Template)
	assert.Equal(t, "waiting", recipe.Status)
	assert.Equal(t, "Running destroy_capsule on cpu.deccd8317431c28552f493a6d4aecf5d.", recipe.StatusDetail)
	assert.Equal(t, "5854017d89d50f424e000002", recipe.AccountID)
	assert.Equal(t, "5854017e89d50f424e000192", recipe.DeploymentID)
}
