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

func TestBroker_FetchServiceInstance(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{"GET", "/deployments", 200, util.Body("../_fixtures/api_get_deployments.json"), nil},
		util.HttpTestCase{"GET", "/deployments/5854017e89d50f424e000192", 200, util.Body("../_fixtures/api_get_deployment.json"), nil},
		util.HttpTestCase{"GET", "/deployments/5854017e89d50f424e000192/recipes", 200, util.Body("../_fixtures/api_get_recipes_for_service_fetch.json"), nil},
		util.HttpTestCase{"GET", "/deployments/5854017e89d50f424e000192/scalings", 200, util.Body("../_fixtures/api_get_scaling.json"), nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, util.Body("../_fixtures/broker_fetch_service_instance.json"), rec.Body.String())
}

func TestBroker_FetchServiceInstance_NotFound(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{"GET", "/deployments", 200, util.Body("../_fixtures/api_get_deployments.json"), nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/v2/service_instances/52551d5f-1350-4f7d-9ddd-710a47ef9b72", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 404, rec.Code)
	assert.Contains(t, rec.Body.String(), `"error": "ServiceInstanceNotFound"`)
	assert.Contains(t, rec.Body.String(), `"description": "The service instance does not exist"`)
}

func TestBroker_FetchServiceInstance_RecipesNotFound(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{"GET", "/deployments", 200, util.Body("../_fixtures/api_get_deployments.json"), nil},
		util.HttpTestCase{"GET", "/deployments/5854017e89d50f424e000192", 200, util.Body("../_fixtures/api_get_deployment.json"), nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 404, rec.Code)
	assert.Contains(t, rec.Body.String(), `"error": "RecipesNotFound"`)
	assert.Contains(t, rec.Body.String(), `"description": "The service instance recipes could not be found"`)
}

func TestBroker_FetchServiceInstance_ScalingsNotFound(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{"GET", "/deployments", 200, util.Body("../_fixtures/api_get_deployments.json"), nil},
		util.HttpTestCase{"GET", "/deployments/5854017e89d50f424e000192", 200, util.Body("../_fixtures/api_get_deployment.json"), nil},
		util.HttpTestCase{"GET", "/deployments/5854017e89d50f424e000192/recipes", 200, util.Body("../_fixtures/api_get_recipes_for_service_fetch.json"), nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 404, rec.Code)
	assert.Contains(t, rec.Body.String(), `"error": "ScalingParametersNotFound"`)
	assert.Contains(t, rec.Body.String(), `"description": "The service instance scaling parameters do not exist"`)
}

func TestBroker_FetchServiceInstance_ConcurrencyError422(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{"GET", "/deployments", 200, util.Body("../_fixtures/api_get_deployments.json"), nil},
		util.HttpTestCase{"GET", "/deployments/5854017e89d50f424e000192", 200, util.Body("../_fixtures/api_get_deployment.json"), nil},
		util.HttpTestCase{"GET", "/deployments/5854017e89d50f424e000192/recipes", 200, util.Body("../_fixtures/api_get_recipes_for_concurrency_error_422.json"), nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 422, rec.Code)
	assert.Contains(t, rec.Body.String(), `"error": "ConcurrencyError"`)
	assert.Contains(t, rec.Body.String(), `"description": "The service instance is being updated"`)
}

func TestBroker_FetchServiceInstance_ConcurrencyError404(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{"GET", "/deployments", 200, util.Body("../_fixtures/api_get_deployments.json"), nil},
		util.HttpTestCase{"GET", "/deployments/5854017e89d50f424e000192", 200, util.Body("../_fixtures/api_get_deployment.json"), nil},
		util.HttpTestCase{"GET", "/deployments/5854017e89d50f424e000192/recipes", 200, util.Body("../_fixtures/api_get_recipes_for_concurrency_error_404.json"), nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 404, rec.Code)
	assert.Contains(t, rec.Body.String(), `"error": "ConcurrencyError"`)
	assert.Contains(t, rec.Body.String(), `"description": "The service instance provisioning is still in progress"`)
}
