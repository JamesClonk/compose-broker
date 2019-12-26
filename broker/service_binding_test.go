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

func TestBroker_BindBinding(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{Method: "GET", Path: "/deployments", Code: 200, Body: util.Body("../_fixtures/api_get_deployments.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192", Code: 200, Body: util.Body("../_fixtures/api_get_deployment_for_service_binding.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192/scalings", Code: 200, Body: util.Body("../_fixtures/api_get_scaling_for_service_binding.json"), Test: nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PUT", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26/service_bindings/deadbeef", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, util.Body("../_fixtures/broker_fetch_service_binding.json"), rec.Body.String())
}

func TestBroker_BindBinding_NotFound(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{Method: "GET", Path: "/deployments", Code: 404, Body: util.Body("../_fixtures/api_get_deployments.json"), Test: nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PUT", "/v2/service_instances/52551d5f-1350-4f7d-9ddd-710a47ef9b72/service_bindings/deadbeef", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 400, rec.Code)
	assert.Contains(t, rec.Body.String(), `"error": "MissingServiceInstance"`)
	assert.Contains(t, rec.Body.String(), `"description": "The service instance does not exist"`)
}

func TestBroker_BindBinding_NoScaling(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{Method: "GET", Path: "/deployments", Code: 200, Body: util.Body("../_fixtures/api_get_deployments.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192", Code: 200, Body: util.Body("../_fixtures/api_get_deployment_for_service_binding.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192/scalings", Code: 404, Body: util.Body("../_fixtures/api_get_scaling.json"), Test: nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PUT", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26/service_bindings/deadbeef", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, util.Body("../_fixtures/broker_fetch_service_binding_without_scaling.json"), rec.Body.String())
}

func TestBroker_FetchBinding(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{Method: "GET", Path: "/deployments", Code: 200, Body: util.Body("../_fixtures/api_get_deployments.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192", Code: 200, Body: util.Body("../_fixtures/api_get_deployment_for_service_binding.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192/scalings", Code: 200, Body: util.Body("../_fixtures/api_get_scaling_for_service_binding.json"), Test: nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26/service_bindings/deadbeef", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, util.Body("../_fixtures/broker_fetch_service_binding.json"), rec.Body.String())
}

func TestBroker_FetchBinding_NotFound(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{Method: "GET", Path: "/deployments", Code: 404, Body: util.Body("../_fixtures/api_get_deployments.json"), Test: nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/v2/service_instances/52551d5f-1350-4f7d-9ddd-710a47ef9b72/service_bindings/deadbeef", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 404, rec.Code)
	assert.Contains(t, rec.Body.String(), `"error": "MissingServiceInstance"`)
	assert.Contains(t, rec.Body.String(), `"description": "The service instance does not exist"`)
}

func TestBroker_FetchBinding_NoScaling(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{Method: "GET", Path: "/deployments", Code: 200, Body: util.Body("../_fixtures/api_get_deployments.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192", Code: 200, Body: util.Body("../_fixtures/api_get_deployment_for_service_binding.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192/scalings", Code: 404, Body: util.Body("../_fixtures/api_get_scaling.json"), Test: nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26/service_bindings/deadbeef", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, util.Body("../_fixtures/broker_fetch_service_binding_without_scaling.json"), rec.Body.String())
}

func TestBroker_UnbindBinding(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{Method: "GET", Path: "/deployments", Code: 200, Body: util.Body("../_fixtures/api_get_deployments.json"), Test: nil},
		util.HttpTestCase{Method: "GET", Path: "/deployments/5854017e89d50f424e000192", Code: 200, Body: util.Body("../_fixtures/api_get_deployment.json"), Test: nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26/service_bindings/deadbeef", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, `{}`, rec.Body.String())
}

func TestBroker_UnbindBinding_NotFound(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{Method: "GET", Path: "/deployments", Code: 404, Body: util.Body("../_fixtures/api_get_deployments.json"), Test: nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", "/v2/service_instances/52551d5f-1350-4f7d-9ddd-710a47ef9b72/service_bindings/deadbeef", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 410, rec.Code)
	assert.Contains(t, rec.Body.String(), `"error": "MissingServiceInstance"`)
	assert.Contains(t, rec.Body.String(), `"description": "The service instance does not exist"`)
}
