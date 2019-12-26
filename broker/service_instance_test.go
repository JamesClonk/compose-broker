package broker

import (
	"bytes"
	"encoding/json"
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

func TestBroker_ProvisionServiceInstance(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{Method: "POST", Path: "/deployments", Code: 202, Body: util.Body("../_fixtures/api_create_deployment_for_service_provisioning.json"), Test: func(body string) {
			assert.Contains(t, body, `"name":"8dcdf609-36c9-4b22-bb16-d97e48c50f26"`)
			assert.Contains(t, body, `"account_id":"586eab527c65836dde5533e8"`)
			assert.Contains(t, body, `"type":"postgresql"`)
			assert.Contains(t, body, `"datacenter":"gce:europe-west1"`)
			assert.Contains(t, body, `"units":1`)
			assert.NotContains(t, body, `version`)
			assert.NotContains(t, body, `cache_mode`)
			assert.Contains(t, body, `"notes":"9b4ee86b-3876-469f-a531-062e71bc5859-d6222855-17c6-448c-885a-e9d931cd221b"`)
		}},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	provisioning := ServiceInstanceProvisioning{
		ServiceID: "9b4ee86b-3876-469f-a531-062e71bc5859",
		PlanID:    "d6222855-17c6-448c-885a-e9d931cd221b",
	}
	data, _ := json.MarshalIndent(provisioning, "", "  ")

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PUT", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26?accepts_incomplete=true", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 202, rec.Code)
	assert.Equal(t, util.Body("../_fixtures/broker_provision_service_instance.json"), rec.Body.String())
}

func TestBroker_ProvisionServiceInstance_WithProvisionParameters(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{Method: "POST", Path: "/deployments", Code: 202, Body: util.Body("../_fixtures/api_create_deployment_for_service_provisioning.json"), Test: func(body string) {
			assert.Contains(t, body, `"name":"8dcdf609-36c9-4b22-bb16-d97e48c50f26"`)
			assert.Contains(t, body, `"account_id":"oracle"`)
			assert.Contains(t, body, `"type":"postgresql"`)
			assert.Contains(t, body, `"datacenter":"solaris:sun"`)
			assert.Contains(t, body, `"units":7`)
			assert.Contains(t, body, `"version":"11.0"`)
			assert.Contains(t, body, `"cache_mode":true`)
			assert.Contains(t, body, `"notes":"9b4ee86b-3876-469f-a531-062e71bc5859-d6222855-17c6-448c-885a-e9d931cd221b"`)
		}},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	provisioning := ServiceInstanceProvisioning{
		ServiceID: "9b4ee86b-3876-469f-a531-062e71bc5859",
		PlanID:    "d6222855-17c6-448c-885a-e9d931cd221b",
	}
	provisioning.Parameters.AccountID = "oracle"
	provisioning.Parameters.Datacenter = "solaris:sun"
	provisioning.Parameters.Units = 7
	provisioning.Parameters.Version = "11.0"
	provisioning.Parameters.CacheMode = true
	data, _ := json.MarshalIndent(provisioning, "", "  ")

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PUT", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26?accepts_incomplete=true", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 202, rec.Code)
	assert.Equal(t, util.Body("../_fixtures/broker_provision_service_instance.json"), rec.Body.String())
}

func TestBroker_ProvisionServiceInstance_WithPlanParameters(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{Method: "POST", Path: "/deployments", Code: 202, Body: util.Body("../_fixtures/api_create_deployment_for_service_provisioning.json"), Test: func(body string) {
			assert.Contains(t, body, `"name":"8dcdf609-36c9-4b22-bb16-d97e48c50f26"`)
			assert.Contains(t, body, `"type":"redis"`)
			assert.Contains(t, body, `"datacenter":"aws:eu-central-1"`)
			assert.Contains(t, body, `"units":2`)
			assert.Contains(t, body, `"version":"4.0.14"`)
			assert.Contains(t, body, `"cache_mode":true`)
			assert.Contains(t, body, `"notes":"e27ea95a-3883-44f2-8ca4-01101f39d50c-ae2bda53-fe15-4335-9422-774aae3e7e32"`)
		}},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	provisioning := ServiceInstanceProvisioning{
		ServiceID: "e27ea95a-3883-44f2-8ca4-01101f39d50c",
		PlanID:    "ae2bda53-fe15-4335-9422-774aae3e7e32",
	}
	data, _ := json.MarshalIndent(provisioning, "", "  ")

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PUT", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26?accepts_incomplete=true", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 202, rec.Code)
	assert.Equal(t, util.Body("../_fixtures/broker_provision_service_instance.json"), rec.Body.String())
}

func TestBroker_ProvisionServiceInstance_AsyncRequired(t *testing.T) {
	r := NewRouter(util.TestConfig(""))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PUT", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 422, rec.Code) // a provisioning request must be async
	assert.Contains(t, rec.Body.String(), `"error": "AsyncRequired"`)
	assert.Contains(t, rec.Body.String(), `"description": "Service instance provisioning requires an asynchronous operation"`)
}

func TestBroker_ProvisionServiceInstance_EmptyBody(t *testing.T) {
	r := NewRouter(util.TestConfig(""))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PUT", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26?accepts_incomplete=true", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 400, rec.Code)
	assert.Contains(t, rec.Body.String(), `"error": "MalformedRequest"`)
	assert.Contains(t, rec.Body.String(), `"description": "Could not read provisioning request"`)
}

func TestBroker_ProvisionServiceInstance_UnknownPlan(t *testing.T) {
	r := NewRouter(util.TestConfig(""))

	provisioning := ServiceInstanceProvisioning{
		ServiceID: "9b4ee86b-3876-469f-a531-062e71bc5859",
		PlanID:    "deadbeef",
	}
	data, _ := json.MarshalIndent(provisioning, "", "  ")

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PUT", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26?accepts_incomplete=true", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 400, rec.Code)
	assert.Contains(t, rec.Body.String(), `"error": "MalformedRequest"`)
	assert.Contains(t, rec.Body.String(), `"description": "Unknown plan_id"`)
}

func TestBroker_ProvisionServiceInstance_UnitsMissing(t *testing.T) {
	r := NewRouter(util.TestConfig(""))

	provisioning := ServiceInstanceProvisioning{
		ServiceID: "9b4ee86b-3876-469f-a531-062e71bc5859",
	}
	data, _ := json.MarshalIndent(provisioning, "", "  ")

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PUT", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26?accepts_incomplete=true", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 400, rec.Code)
	assert.Contains(t, rec.Body.String(), `"error": "MissingParameters"`)
	assert.Contains(t, rec.Body.String(), `"description": "Units parameter is missing for service instance provisioning"`)
}

func TestBroker_ProvisionServiceInstance_AccountsNotReadable(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{Method: "GET", Path: "/accounts", Code: 404, Body: util.Body("../_fixtures/api_get_accounts_invalid_file.json"), Test: nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()

	config := util.TestConfig(apiServer.URL)
	config.API.DefaultAccountID = "" // clear
	r := NewRouter(config)

	provisioning := ServiceInstanceProvisioning{
		ServiceID: "9b4ee86b-3876-469f-a531-062e71bc5859",
		PlanID:    "d6222855-17c6-448c-885a-e9d931cd221b",
	}
	data, _ := json.MarshalIndent(provisioning, "", "  ")

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PUT", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26?accepts_incomplete=true", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 409, rec.Code)
	assert.Contains(t, rec.Body.String(), `"error": "UnknownError"`)
	assert.Contains(t, rec.Body.String(), `"description": "Could not read Compose.io accounts"`)
}

func TestBroker_ProvisionServiceInstance_AccountIDFromAPI(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{Method: "GET", Path: "/accounts", Code: 200, Body: util.Body("../_fixtures/api_get_accounts_for_service_provision.json"), Test: nil},
		util.HttpTestCase{Method: "POST", Path: "/deployments", Code: 202, Body: util.Body("../_fixtures/api_create_deployment_for_service_provisioning.json"), Test: func(body string) {
			assert.Contains(t, body, `"name":"8dcdf609-36c9-4b22-bb16-d97e48c50f26"`)
			assert.Contains(t, body, `"account_id":"22de8c5fbdc3d1f777750492"`)
			assert.Contains(t, body, `"type":"postgresql"`)
			assert.Contains(t, body, `"datacenter":"gce:europe-west1"`)
			assert.Contains(t, body, `"units":1`)
			assert.NotContains(t, body, `version`)
			assert.NotContains(t, body, `cache_mode`)
			assert.Contains(t, body, `"notes":"9b4ee86b-3876-469f-a531-062e71bc5859-d6222855-17c6-448c-885a-e9d931cd221b"`)
		}},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()

	config := util.TestConfig(apiServer.URL)
	config.API.DefaultAccountID = "" // clear
	r := NewRouter(config)

	provisioning := ServiceInstanceProvisioning{
		ServiceID: "9b4ee86b-3876-469f-a531-062e71bc5859",
		PlanID:    "d6222855-17c6-448c-885a-e9d931cd221b",
	}
	data, _ := json.MarshalIndent(provisioning, "", "  ")

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PUT", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26?accepts_incomplete=true", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 202, rec.Code)
	assert.Equal(t, util.Body("../_fixtures/broker_provision_service_instance.json"), rec.Body.String())
}

func TestBroker_ProvisionServiceInstance_AccountIDMissing(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{Method: "GET", Path: "/accounts", Code: 200, Body: util.Body("../_fixtures/api_get_accounts_empty.json"), Test: nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()

	config := util.TestConfig(apiServer.URL)
	config.API.DefaultAccountID = "" // clear
	r := NewRouter(config)

	provisioning := ServiceInstanceProvisioning{
		ServiceID: "9b4ee86b-3876-469f-a531-062e71bc5859",
		PlanID:    "d6222855-17c6-448c-885a-e9d931cd221b",
	}
	data, _ := json.MarshalIndent(provisioning, "", "  ")

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PUT", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26?accepts_incomplete=true", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 400, rec.Code)
	assert.Contains(t, rec.Body.String(), `"error": "MissingParameters"`)
	assert.Contains(t, rec.Body.String(), `"description": "AccountID is missing for service instance provisioning"`)
}

func TestBroker_ProvisionServiceInstance_AlreadyExistsWithSameScaling(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{Method: "GET", Path: "/deployments", Code: 200, Body: util.Body("../_fixtures/api_get_deployments.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192", Code: 200, Body: util.Body("../_fixtures/api_get_deployment.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192/recipes", Code: 200, Body: util.Body("../_fixtures/api_get_recipes_for_service_provision_200.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192/scalings", Code: 200, Body: util.Body("../_fixtures/api_get_scaling_for_service_provision.json"), Test: nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	provisioning := ServiceInstanceProvisioning{
		ServiceID: "9b4ee86b-3876-469f-a531-062e71bc5859",
		PlanID:    "d6222855-17c6-448c-885a-e9d931cd221b",
	}
	provisioning.Parameters.Units = 5
	data, _ := json.MarshalIndent(provisioning, "", "  ")

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PUT", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26?accepts_incomplete=true", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, util.Body("../_fixtures/broker_provision_service_instance_existing.json"), rec.Body.String())
}

func TestBroker_ProvisionServiceInstance_AlreadyExistsWithOngoingRecipe(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{Method: "GET", Path: "/deployments", Code: 200, Body: util.Body("../_fixtures/api_get_deployments.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192", Code: 200, Body: util.Body("../_fixtures/api_get_deployment.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192/recipes", Code: 200, Body: util.Body("../_fixtures/api_get_recipes_for_service_provision_202.json"), Test: nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	provisioning := ServiceInstanceProvisioning{
		ServiceID: "9b4ee86b-3876-469f-a531-062e71bc5859",
		PlanID:    "d6222855-17c6-448c-885a-e9d931cd221b",
	}
	data, _ := json.MarshalIndent(provisioning, "", "  ")

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PUT", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26?accepts_incomplete=true", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 202, rec.Code)
	assert.Equal(t, util.Body("../_fixtures/broker_provision_service_instance_existing.json"), rec.Body.String())
}

func TestBroker_ProvisionServiceInstance_AlreadyExistsButNoRecipes(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{Method: "GET", Path: "/deployments", Code: 200, Body: util.Body("../_fixtures/api_get_deployments.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192", Code: 200, Body: util.Body("../_fixtures/api_get_deployment.json"), Test: nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	provisioning := ServiceInstanceProvisioning{
		ServiceID: "9b4ee86b-3876-469f-a531-062e71bc5859",
		PlanID:    "d6222855-17c6-448c-885a-e9d931cd221b",
	}
	data, _ := json.MarshalIndent(provisioning, "", "  ")

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PUT", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26?accepts_incomplete=true", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 409, rec.Code)
	assert.Contains(t, rec.Body.String(), `"error": "UnknownError"`)
	assert.Contains(t, rec.Body.String(), `"description": "Could not create service instance"`)
}

func TestBroker_ProvisionServiceInstance_Error(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{Method: "POST", Path: "/deployments", Code: 418, Body: util.Body("../_fixtures/api_create_deployment_for_service_provisioning.json"), Test: nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	provisioning := ServiceInstanceProvisioning{
		ServiceID: "9b4ee86b-3876-469f-a531-062e71bc5859",
		PlanID:    "d6222855-17c6-448c-885a-e9d931cd221b",
	}
	data, _ := json.MarshalIndent(provisioning, "", "  ")

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PUT", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26?accepts_incomplete=true", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 500, rec.Code)
	assert.Contains(t, rec.Body.String(), `"error": "UnknownError"`)
	assert.Contains(t, rec.Body.String(), `"description": "Could not create service instance"`)
}

func TestBroker_ProvisionServiceInstance_ImmediateCreation(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{Method: "POST", Path: "/deployments", Code: 202, Body: util.Body("../_fixtures/api_create_deployment_for_service_provisioning.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/recipes/59a6b3a5f32fb6001001ae6b", Code: 200, Body: util.Body("../_fixtures/api_get_recipe_for_immediate_service_provision.json"), Test: nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	provisioning := ServiceInstanceProvisioning{
		ServiceID: "9b4ee86b-3876-469f-a531-062e71bc5859",
		PlanID:    "d6222855-17c6-448c-885a-e9d931cd221b",
	}
	data, _ := json.MarshalIndent(provisioning, "", "  ")

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PUT", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26?accepts_incomplete=true", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 201, rec.Code) // provisioning could be fast and be immediately done
	assert.Equal(t, "{}", rec.Body.String())
}

func TestBroker_ProvisionServiceInstance_ImmediateFailure(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{Method: "POST", Path: "/deployments", Code: 202, Body: util.Body("../_fixtures/api_create_deployment_for_service_provisioning.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/recipes/59a6b3a5f32fb6001001ae6b", Code: 200, Body: util.Body("../_fixtures/api_get_recipe_for_immediate_service_provision_failure.json"), Test: nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	provisioning := ServiceInstanceProvisioning{
		ServiceID: "9b4ee86b-3876-469f-a531-062e71bc5859",
		PlanID:    "d6222855-17c6-448c-885a-e9d931cd221b",
	}
	data, _ := json.MarshalIndent(provisioning, "", "  ")

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PUT", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26?accepts_incomplete=true", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 400, rec.Code)
	assert.Contains(t, rec.Body.String(), `"error": "ProvisionFailure"`)
	assert.Contains(t, rec.Body.String(), `"description": "Could not create service instance"`)
}

func TestBroker_LastOperationServiceInstance(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{Method: "GET", Path: "/deployments", Code: 200, Body: util.Body("../_fixtures/api_get_deployments.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192", Code: 200, Body: util.Body("../_fixtures/api_get_deployment.json"), Test: nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26/last_operation", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, util.Body("../_fixtures/broker_last_operation_on_service_instance_no_recipe.json"), rec.Body.String())
}

func TestBroker_LastOperationServiceInstance_Succeeded(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{Method: "GET", Path: "/deployments", Code: 200, Body: util.Body("../_fixtures/api_get_deployments.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192", Code: 200, Body: util.Body("../_fixtures/api_get_deployment.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192/recipes", Code: 200, Body: util.Body("../_fixtures/api_get_recipes_for_last_operation_succeeded.json"), Test: nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26/last_operation", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, util.Body("../_fixtures/broker_last_operation_on_service_instance_succeeded.json"), rec.Body.String())
}

func TestBroker_LastOperationServiceInstance_Failed(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{Method: "GET", Path: "/deployments", Code: 200, Body: util.Body("../_fixtures/api_get_deployments.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192", Code: 200, Body: util.Body("../_fixtures/api_get_deployment.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192/recipes", Code: 200, Body: util.Body("../_fixtures/api_get_recipes_for_last_operation_failed.json"), Test: nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26/last_operation", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, util.Body("../_fixtures/broker_last_operation_on_service_instance_failed.json"), rec.Body.String())
}

func TestBroker_LastOperationServiceInstance_InProgress(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{Method: "GET", Path: "/deployments", Code: 200, Body: util.Body("../_fixtures/api_get_deployments.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192", Code: 200, Body: util.Body("../_fixtures/api_get_deployment.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192/recipes", Code: 200, Body: util.Body("../_fixtures/api_get_recipes_for_service_provision_in_progress.json"), Test: nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26/last_operation", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, util.Body("../_fixtures/broker_last_operation_on_service_instance_in_progress.json"), rec.Body.String())
}

func TestBroker_LastOperationServiceInstance_NotFound(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{Method: "GET", Path: "/deployments", Code: 404, Body: util.Body("../_fixtures/api_get_deployments.json"), Test: nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26/last_operation", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 410, rec.Code)
	assert.Contains(t, rec.Body.String(), `"error": "MissingServiceInstance"`)
	assert.Contains(t, rec.Body.String(), `"description": "The service instance does not exist"`)
}

func TestBroker_FetchServiceInstance(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{Method: "GET", Path: "/deployments", Code: 200, Body: util.Body("../_fixtures/api_get_deployments.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192", Code: 200, Body: util.Body("../_fixtures/api_get_deployment.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192/recipes", Code: 200, Body: util.Body("../_fixtures/api_get_recipes_for_service_fetch.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192/scalings", Code: 200, Body: util.Body("../_fixtures/api_get_scaling.json"), Test: nil},
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
		util.HttpTestCase{Method: "GET", Path: "/deployments", Code: 404, Body: util.Body("../_fixtures/api_get_deployments.json"), Test: nil},
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
	assert.Contains(t, rec.Body.String(), `"error": "MissingServiceInstance"`)
	assert.Contains(t, rec.Body.String(), `"description": "The service instance does not exist"`)
}

func TestBroker_FetchServiceInstance_RecipesNotFound(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{Method: "GET", Path: "/deployments", Code: 200, Body: util.Body("../_fixtures/api_get_deployments.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192", Code: 200, Body: util.Body("../_fixtures/api_get_deployment.json"), Test: nil},
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
	assert.Contains(t, rec.Body.String(), `"error": "MissingRecipes"`)
	assert.Contains(t, rec.Body.String(), `"description": "The service instance recipes could not be found"`)
}

func TestBroker_FetchServiceInstance_ScalingsNotFound(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{Method: "GET", Path: "/deployments", Code: 200, Body: util.Body("../_fixtures/api_get_deployments.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192", Code: 200, Body: util.Body("../_fixtures/api_get_deployment.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192/recipes", Code: 200, Body: util.Body("../_fixtures/api_get_recipes_for_service_fetch.json"), Test: nil},
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
	assert.Contains(t, rec.Body.String(), `"error": "MissingScalingParameters"`)
	assert.Contains(t, rec.Body.String(), `"description": "The service instance scaling parameters do not exist"`)
}

func TestBroker_FetchServiceInstance_ConcurrencyError422(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{Method: "GET", Path: "/deployments", Code: 200, Body: util.Body("../_fixtures/api_get_deployments.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192", Code: 200, Body: util.Body("../_fixtures/api_get_deployment.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192/recipes", Code: 200, Body: util.Body("../_fixtures/api_get_recipes_for_concurrency_error_422.json"), Test: nil},
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
		util.HttpTestCase{Method: "GET", Path: "/deployments", Code: 200, Body: util.Body("../_fixtures/api_get_deployments.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192", Code: 200, Body: util.Body("../_fixtures/api_get_deployment.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192/recipes", Code: 200, Body: util.Body("../_fixtures/api_get_recipes_for_concurrency_error_404.json"), Test: nil},
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

func TestBroker_UpdateServiceInstance(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{Method: "GET", Path: "/deployments", Code: 200, Body: util.Body("../_fixtures/api_get_deployments.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192", Code: 200, Body: util.Body("../_fixtures/api_get_deployment.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192/scalings", Code: 200, Body: util.Body("../_fixtures/api_get_scaling.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192/recipes", Code: 200, Body: util.Body("../_fixtures/api_get_recipes_for_service_update.json"), Test: nil},
		util.HttpTestCase{Method: "POST", Path: "/deployments/5854017e89d50f424e000192/scalings", Code: 200, Body: util.Body("../_fixtures/api_update_scaling.json"), Test: func(body string) {
			assert.Contains(t, body, `{"deployment":{"units":1}}`)
		}},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	update := ServiceInstanceUpdate{
		ServiceID: "9b4ee86b-3876-469f-a531-062e71bc5859",
		PlanID:    "d6222855-17c6-448c-885a-e9d931cd221b",
	}
	data, _ := json.MarshalIndent(update, "", "  ")

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PATCH", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26?accepts_incomplete=true", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 202, rec.Code)
	assert.Equal(t, util.Body("../_fixtures/broker_update_service_instance.json"), rec.Body.String())
}

func TestBroker_UpdateServiceInstance_NoUpdate(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{Method: "GET", Path: "/deployments", Code: 200, Body: util.Body("../_fixtures/api_get_deployments.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192", Code: 200, Body: util.Body("../_fixtures/api_get_deployment.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192/scalings", Code: 200, Body: util.Body("../_fixtures/api_get_scaling.json"), Test: nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	update := ServiceInstanceUpdate{
		ServiceID: "9b4ee86b-3876-469f-a531-062e71bc5859",
	}
	update.Parameters.Units = 4
	data, _ := json.MarshalIndent(update, "", "  ")

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PATCH", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26?accepts_incomplete=true", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, "{}", rec.Body.String())
}

func TestBroker_UpdateServiceInstance_ImmediateUpdate(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{Method: "GET", Path: "/deployments", Code: 200, Body: util.Body("../_fixtures/api_get_deployments.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192", Code: 200, Body: util.Body("../_fixtures/api_get_deployment.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192/scalings", Code: 200, Body: util.Body("../_fixtures/api_get_scaling.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192/recipes", Code: 200, Body: util.Body("../_fixtures/api_get_recipes_for_service_update.json"), Test: nil},
		util.HttpTestCase{Method: "POST", Path: "/deployments/5854017e89d50f424e000192/scalings", Code: 200, Body: util.Body("../_fixtures/api_update_scaling_for_service_update.json"), Test: func(body string) {
			assert.Contains(t, body, `{"deployment":{"units":1}}`)
		}},
		util.HttpTestCase{Method: "GET", Path: "/recipes/5821fd28a4b549d06e39886d", Code: 200, Body: util.Body("../_fixtures/api_get_recipe_for_immediate_service_update.json"), Test: nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	update := ServiceInstanceUpdate{
		ServiceID: "9b4ee86b-3876-469f-a531-062e71bc5859",
		PlanID:    "d6222855-17c6-448c-885a-e9d931cd221b",
	}
	data, _ := json.MarshalIndent(update, "", "  ")

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PATCH", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26?accepts_incomplete=true", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, "{}", rec.Body.String())
}

func TestBroker_UpdateServiceInstance_ImmediateFailure(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{Method: "GET", Path: "/deployments", Code: 200, Body: util.Body("../_fixtures/api_get_deployments.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192", Code: 200, Body: util.Body("../_fixtures/api_get_deployment.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192/scalings", Code: 200, Body: util.Body("../_fixtures/api_get_scaling.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192/recipes", Code: 200, Body: util.Body("../_fixtures/api_get_recipes_for_service_update.json"), Test: nil},
		util.HttpTestCase{Method: "POST", Path: "/deployments/5854017e89d50f424e000192/scalings", Code: 200, Body: util.Body("../_fixtures/api_update_scaling_for_service_update.json"), Test: func(body string) {
			assert.Contains(t, body, `{"deployment":{"units":1}}`)
		}},
		util.HttpTestCase{Method: "GET", Path: "/recipes/5821fd28a4b549d06e39886d", Code: 200, Body: util.Body("../_fixtures/api_get_recipe_for_immediate_service_update_failure.json"), Test: nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	update := ServiceInstanceUpdate{
		ServiceID: "9b4ee86b-3876-469f-a531-062e71bc5859",
		PlanID:    "d6222855-17c6-448c-885a-e9d931cd221b",
	}
	data, _ := json.MarshalIndent(update, "", "  ")

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PATCH", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26?accepts_incomplete=true", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 409, rec.Code)
	assert.Contains(t, rec.Body.String(), `"error": "UpdateFailure"`)
	assert.Contains(t, rec.Body.String(), `"description": "Could not update service instance"`)
}

func TestBroker_UpdateServiceInstance_WithUnitsParameter(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{Method: "GET", Path: "/deployments", Code: 200, Body: util.Body("../_fixtures/api_get_deployments.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192", Code: 200, Body: util.Body("../_fixtures/api_get_deployment.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192/scalings", Code: 200, Body: util.Body("../_fixtures/api_get_scaling.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192/recipes", Code: 200, Body: util.Body("../_fixtures/api_get_recipes_for_service_update.json"), Test: nil},
		util.HttpTestCase{Method: "POST", Path: "/deployments/5854017e89d50f424e000192/scalings", Code: 200, Body: util.Body("../_fixtures/api_update_scaling.json"), Test: func(body string) {
			assert.Contains(t, body, `{"deployment":{"units":5}}`)
		}},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	update := ServiceInstanceUpdate{
		ServiceID: "9b4ee86b-3876-469f-a531-062e71bc5859",
		PlanID:    "",
	}
	update.Parameters.Units = 5
	data, _ := json.MarshalIndent(update, "", "  ")

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PATCH", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26?accepts_incomplete=true", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 202, rec.Code)
	assert.Equal(t, util.Body("../_fixtures/broker_update_service_instance.json"), rec.Body.String())
}

func TestBroker_UpdateServiceInstance_AsyncRequired(t *testing.T) {
	r := NewRouter(util.TestConfig(""))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PATCH", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 422, rec.Code) // an update request must be async
	assert.Contains(t, rec.Body.String(), `"error": "AsyncRequired"`)
	assert.Contains(t, rec.Body.String(), `"description": "Service instance updating requires an asynchronous operation"`)
}

func TestBroker_UpdateServiceInstance_EmptyBody(t *testing.T) {
	r := NewRouter(util.TestConfig(""))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PATCH", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26?accepts_incomplete=true", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 400, rec.Code)
	assert.Contains(t, rec.Body.String(), `"error": "MalformedRequest"`)
	assert.Contains(t, rec.Body.String(), `"description": "Could not read update request"`)
}

func TestBroker_UpdateServiceInstance_UnknownPlan(t *testing.T) {
	r := NewRouter(util.TestConfig(""))

	update := ServiceInstanceUpdate{
		ServiceID: "9b4ee86b-3876-469f-a531-062e71bc5859",
		PlanID:    "deadbeef",
	}
	data, _ := json.MarshalIndent(update, "", "  ")

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PATCH", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26?accepts_incomplete=true", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 400, rec.Code)
	assert.Contains(t, rec.Body.String(), `"error": "MalformedRequest"`)
	assert.Contains(t, rec.Body.String(), `"description": "Unknown plan_id"`)
}

func TestBroker_UpdateServiceInstance_UnitsMissing(t *testing.T) {
	r := NewRouter(util.TestConfig(""))

	update := ServiceInstanceUpdate{
		ServiceID: "9b4ee86b-3876-469f-a531-062e71bc5859",
	}
	data, _ := json.MarshalIndent(update, "", "  ")

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PATCH", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26?accepts_incomplete=true", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 400, rec.Code)
	assert.Contains(t, rec.Body.String(), `"error": "MissingParameters"`)
	assert.Contains(t, rec.Body.String(), `"description": "Units parameter is missing for service instance update"`)
}

func TestBroker_UpdateServiceInstance_NotFound(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{Method: "GET", Path: "/deployments", Code: 404, Body: util.Body("../_fixtures/api_get_deployments.json"), Test: nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	update := ServiceInstanceUpdate{
		ServiceID: "9b4ee86b-3876-469f-a531-062e71bc5859",
		PlanID:    "d6222855-17c6-448c-885a-e9d931cd221b",
	}
	data, _ := json.MarshalIndent(update, "", "  ")

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PATCH", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26?accepts_incomplete=true", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 404, rec.Code)
	assert.Contains(t, rec.Body.String(), `"error": "ServiceInstanceNotFound"`)
	assert.Contains(t, rec.Body.String(), `"description": "The service instance does not exist"`)
}

func TestBroker_UpdateServiceInstance_ScalingNotFound(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{Method: "GET", Path: "/deployments", Code: 200, Body: util.Body("../_fixtures/api_get_deployments.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192", Code: 200, Body: util.Body("../_fixtures/api_get_deployment.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192/scalings", Code: 404, Body: util.Body("../_fixtures/api_get_scaling.json"), Test: nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	update := ServiceInstanceUpdate{
		ServiceID: "9b4ee86b-3876-469f-a531-062e71bc5859",
		PlanID:    "d6222855-17c6-448c-885a-e9d931cd221b",
	}
	data, _ := json.MarshalIndent(update, "", "  ")

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PATCH", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26?accepts_incomplete=true", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 409, rec.Code)
	assert.Contains(t, rec.Body.String(), `"error": "UnknownError"`)
	assert.Contains(t, rec.Body.String(), `"description": "Could not read service instance scaling"`)
}

func TestBroker_UpdateServiceInstance_ConcurrencyError(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{Method: "GET", Path: "/deployments", Code: 200, Body: util.Body("../_fixtures/api_get_deployments.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192", Code: 200, Body: util.Body("../_fixtures/api_get_deployment.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192/scalings", Code: 200, Body: util.Body("../_fixtures/api_get_scaling.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192/recipes", Code: 200, Body: util.Body("../_fixtures/api_get_recipes_for_concurrency_error_422.json"), Test: nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	update := ServiceInstanceUpdate{
		ServiceID: "9b4ee86b-3876-469f-a531-062e71bc5859",
		PlanID:    "d6222855-17c6-448c-885a-e9d931cd221b",
	}
	data, _ := json.MarshalIndent(update, "", "  ")

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PATCH", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26?accepts_incomplete=true", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 422, rec.Code)
	assert.Contains(t, rec.Body.String(), `"error": "ConcurrencyError"`)
	assert.Contains(t, rec.Body.String(), `"description": "The service instance is currently being updated"`)
}

func TestBroker_UpdateServiceInstance_Error(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{Method: "GET", Path: "/deployments", Code: 200, Body: util.Body("../_fixtures/api_get_deployments.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192", Code: 200, Body: util.Body("../_fixtures/api_get_deployment.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192/scalings", Code: 200, Body: util.Body("../_fixtures/api_get_scaling.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192/recipes", Code: 200, Body: util.Body("../_fixtures/api_get_recipes_for_service_update.json"), Test: nil},
		util.HttpTestCase{Method: "POST", Path: "/deployments/5854017e89d50f424e000192/scalings", Code: 500, Body: util.Body("../_fixtures/api_update_scaling.json"), Test: func(body string) {
			assert.Contains(t, body, `{"deployment":{"units":1}}`)
		}},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	update := ServiceInstanceUpdate{
		ServiceID: "9b4ee86b-3876-469f-a531-062e71bc5859",
		PlanID:    "d6222855-17c6-448c-885a-e9d931cd221b",
	}
	data, _ := json.MarshalIndent(update, "", "  ")

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PATCH", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26?accepts_incomplete=true", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 409, rec.Code)
	assert.Contains(t, rec.Body.String(), `"error": "UnknownError"`)
	assert.Contains(t, rec.Body.String(), `"description": "Could not update service instance"`)
}

func TestBroker_DeprovisionServiceInstance(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{Method: "GET", Path: "/deployments", Code: 200, Body: util.Body("../_fixtures/api_get_deployments.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192", Code: 200, Body: util.Body("../_fixtures/api_get_deployment.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192/recipes", Code: 200, Body: util.Body("../_fixtures/api_get_recipes_for_service_deprovision.json"), Test: nil},
		util.HttpTestCase{Method: "DELETE", Path: "/deployments/5854017e89d50f424e000192", Code: 202, Body: util.Body("../_fixtures/api_delete_deployment_for_service_deprovision.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/recipes/5821fd28a4b549d06e39886d", Code: 200, Body: util.Body("../_fixtures/api_get_recipe_for_service_deprovision.json"), Test: nil},
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
		util.HttpTestCase{Method: "GET", Path: "/deployments", Code: 200, Body: util.Body("../_fixtures/api_get_deployments.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192", Code: 200, Body: util.Body("../_fixtures/api_get_deployment.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192/recipes", Code: 200, Body: util.Body("../_fixtures/api_get_recipes_for_service_deprovision.json"), Test: nil},
		util.HttpTestCase{Method: "DELETE", Path: "/deployments/5854017e89d50f424e000192", Code: 202, Body: util.Body("../_fixtures/api_delete_deployment_for_service_deprovision.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/recipes/5821fd28a4b549d06e39886d", Code: 200, Body: util.Body("../_fixtures/api_get_recipe_for_immediate_service_deprovision.json"), Test: nil},
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
		util.HttpTestCase{Method: "GET", Path: "/deployments", Code: 200, Body: util.Body("../_fixtures/api_get_deployments.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192", Code: 200, Body: util.Body("../_fixtures/api_get_deployment.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192/recipes", Code: 200, Body: util.Body("../_fixtures/api_get_recipes_for_service_deprovision.json"), Test: nil},
		util.HttpTestCase{Method: "DELETE", Path: "/deployments/5854017e89d50f424e000192", Code: 202, Body: util.Body("../_fixtures/api_delete_deployment_for_service_deprovision.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/recipes/5821fd28a4b549d06e39886d", Code: 200, Body: util.Body("../_fixtures/api_get_recipe_for_immediate_service_deprovision_failure.json"), Test: nil},
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
		util.HttpTestCase{Method: "GET", Path: "/deployments", Code: 404, Body: util.Body("../_fixtures/api_get_deployments.json"), Test: nil},
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
		util.HttpTestCase{Method: "GET", Path: "/deployments", Code: 200, Body: util.Body("../_fixtures/api_get_deployments.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192", Code: 200, Body: util.Body("../_fixtures/api_get_deployment.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192/recipes", Code: 200, Body: util.Body("../_fixtures/api_get_recipes_for_concurrency_error_422.json"), Test: nil},
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
	assert.Contains(t, rec.Body.String(), `"description": "The service instance is currently being updated"`)
}

func TestBroker_DeprovisionServiceInstance_Error(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{Method: "GET", Path: "/deployments", Code: 200, Body: util.Body("../_fixtures/api_get_deployments.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192", Code: 200, Body: util.Body("../_fixtures/api_get_deployment.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192/recipes", Code: 200, Body: util.Body("../_fixtures/api_get_recipes_for_service_deprovision.json"), Test: nil},
		util.HttpTestCase{Method: "DELETE", Path: "/deployments/5854017e89d50f424e000192", Code: 500, Body: util.Body("../_fixtures/api_delete_deployment.json"), Test: nil},
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
