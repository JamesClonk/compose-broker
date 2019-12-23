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
		util.HttpTestCase{"GET", "/deployments", 404, util.Body("../_fixtures/api_get_deployments.json"), nil},
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

func TestBroker_DeprovisionServiceInstance(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{"GET", "/deployments", 200, util.Body("../_fixtures/api_get_deployments.json"), nil},
		util.HttpTestCase{"GET", "/deployments/5854017e89d50f424e000192", 200, util.Body("../_fixtures/api_get_deployment.json"), nil},
		util.HttpTestCase{"GET", "/deployments/5854017e89d50f424e000192/recipes", 200, util.Body("../_fixtures/api_get_recipes_for_service_deprovision.json"), nil},
		util.HttpTestCase{"DELETE", "/deployments/5854017e89d50f424e000192", 202, util.Body("../_fixtures/api_delete_deployment_for_service_deprovision.json"), nil},
		util.HttpTestCase{"GET", "/recipes/5821fd28a4b549d06e39886d", 200, util.Body("../_fixtures/api_get_recipe_for_service_deprovision.json"), nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26?accepts_incomplete=true", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 202, rec.Code) // a normal deprovisioning should be async
	assert.Equal(t, "{}", rec.Body.String())
}

func TestBroker_DeprovisionServiceInstance_ImmediateDeletion(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{"GET", "/deployments", 200, util.Body("../_fixtures/api_get_deployments.json"), nil},
		util.HttpTestCase{"GET", "/deployments/5854017e89d50f424e000192", 200, util.Body("../_fixtures/api_get_deployment.json"), nil},
		util.HttpTestCase{"GET", "/deployments/5854017e89d50f424e000192/recipes", 200, util.Body("../_fixtures/api_get_recipes_for_service_deprovision.json"), nil},
		util.HttpTestCase{"DELETE", "/deployments/5854017e89d50f424e000192", 202, util.Body("../_fixtures/api_delete_deployment_for_service_deprovision.json"), nil},
		util.HttpTestCase{"GET", "/recipes/5821fd28a4b549d06e39886d", 200, util.Body("../_fixtures/api_get_recipe_for_immediate_service_deprovision.json"), nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26?accepts_incomplete=true", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code) // deprovisioning could be fast and be immediately done
	assert.Equal(t, "{}", rec.Body.String())
}

func TestBroker_DeprovisionServiceInstance_ImmediateFailure(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{"GET", "/deployments", 200, util.Body("../_fixtures/api_get_deployments.json"), nil},
		util.HttpTestCase{"GET", "/deployments/5854017e89d50f424e000192", 200, util.Body("../_fixtures/api_get_deployment.json"), nil},
		util.HttpTestCase{"GET", "/deployments/5854017e89d50f424e000192/recipes", 200, util.Body("../_fixtures/api_get_recipes_for_service_deprovision.json"), nil},
		util.HttpTestCase{"DELETE", "/deployments/5854017e89d50f424e000192", 202, util.Body("../_fixtures/api_delete_deployment_for_service_deprovision.json"), nil},
		util.HttpTestCase{"GET", "/recipes/5821fd28a4b549d06e39886d", 200, util.Body("../_fixtures/api_get_recipe_for_immediate_service_deprovision_failure.json"), nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26?accepts_incomplete=true", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 500, rec.Code)
	assert.Contains(t, rec.Body.String(), `"error": "DeprovisionFailure"`)
	assert.Contains(t, rec.Body.String(), `"description": "Could not delete service instance"`)
}

func TestBroker_DeprovisionServiceInstance_AsyncRequired(t *testing.T) {
	r := NewRouter(util.TestConfig(""))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 422, rec.Code) // a deprovisioning request must be async
	assert.Contains(t, rec.Body.String(), `"error": "AsyncRequired"`)
	assert.Contains(t, rec.Body.String(), `"description": "Service instance deprovisioning requires an asynchronous operation"`)
}

func TestBroker_DeprovisionServiceInstance_NotFound(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{"GET", "/deployments", 404, util.Body("../_fixtures/api_get_deployments.json"), nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26?accepts_incomplete=true", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 410, rec.Code)
	assert.Contains(t, rec.Body.String(), `"error": "MissingServiceInstance"`)
	assert.Contains(t, rec.Body.String(), `"description": "The service instance does not exist"`)
}

func TestBroker_DeprovisionServiceInstance_ConcurrencyError(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{"GET", "/deployments", 200, util.Body("../_fixtures/api_get_deployments.json"), nil},
		util.HttpTestCase{"GET", "/deployments/5854017e89d50f424e000192", 200, util.Body("../_fixtures/api_get_deployment.json"), nil},
		util.HttpTestCase{"GET", "/deployments/5854017e89d50f424e000192/recipes", 200, util.Body("../_fixtures/api_get_recipes_for_concurrency_error_422.json"), nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26?accepts_incomplete=true", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 422, rec.Code)
	assert.Contains(t, rec.Body.String(), `"error": "ConcurrencyError"`)
	assert.Contains(t, rec.Body.String(), `"description": "The service instance is being updated"`)
}

func TestBroker_DeprovisionServiceInstance_Error(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{"GET", "/deployments", 200, util.Body("../_fixtures/api_get_deployments.json"), nil},
		util.HttpTestCase{"GET", "/deployments/5854017e89d50f424e000192", 200, util.Body("../_fixtures/api_get_deployment.json"), nil},
		util.HttpTestCase{"GET", "/deployments/5854017e89d50f424e000192/recipes", 200, util.Body("../_fixtures/api_get_recipes_for_service_deprovision.json"), nil},
		util.HttpTestCase{"DELETE", "/deployments/5854017e89d50f424e000192", 500, util.Body("../_fixtures/api_delete_deployment.json"), nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26?accepts_incomplete=true", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 500, rec.Code)
	assert.Contains(t, rec.Body.String(), `"error": "UnknownError"`)
	assert.Contains(t, rec.Body.String(), `"description": "Could not delete service instance"`)
}
