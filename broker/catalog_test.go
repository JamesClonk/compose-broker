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

func TestBroker_Catalog(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{"GET", "/databases", 200, util.Body("../_fixtures/api_get_databases.json"), nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/v2/catalog", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
	assert.Contains(t, rec.Body.String(), `"id": "b0d27854-06e0-426e-9e8d-79f4c53078c7"`)
	assert.Contains(t, rec.Body.String(), `"id": "355ef4a4-08f5-4764-b4ed-8353812b6963"`)
	assert.Contains(t, rec.Body.String(), `"displayName": "PostgreSQL"`)
	assert.Contains(t, rec.Body.String(), `"name": "rethink"`)
	assert.Contains(t, rec.Body.String(), `"documentationUrl": "https://compose.com/databases/scylladb"`)
	assert.Contains(t, rec.Body.String(), `"imageUrl": "https://compose.com/assets/icd-icons/etcd-9bf4cedacb096e58868085dc0b91b8c3cc3c1f0b3f8be05d1bf7ea5b5e9b6697.svg"`)
	assert.Contains(t, rec.Body.String(), `"2 GB Storage"`)
	assert.Contains(t, rec.Body.String(), `"512 MB RAM"`)
	assert.Contains(t, rec.Body.String(), `"datacenter": "aws:eu-central-1"`)
	assert.Contains(t, rec.Body.String(), `"longDescription": "Deploy RabbitMQ on AWS, GCP, or IBM Cloud in minutes. Fully managed, highly-available and production ready."`)
	assert.Equal(t, util.Body("../_fixtures/broker_catalog.json"), rec.Body.String())
}

func TestBroker_Catalog_Trimmed(t *testing.T) {
	test := []util.HttpTestCase{
		util.HttpTestCase{"GET", "/databases", 200, util.Body("../_fixtures/api_get_databases_trimmed.json"), nil},
	}
	apiServer := util.TestServer("deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/v2/catalog", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
	assert.Contains(t, rec.Body.String(), `"id": "b0d27854-06e0-426e-9e8d-79f4c53078c7"`)
	assert.Contains(t, rec.Body.String(), `"id": "355ef4a4-08f5-4764-b4ed-8353812b6963"`)
	assert.Contains(t, rec.Body.String(), `"displayName": "PostgreSQL"`)
	assert.Contains(t, rec.Body.String(), `"name": "rethink"`)
	assert.NotContains(t, rec.Body.String(), `"etcd"`)
	assert.NotContains(t, rec.Body.String(), `"scylla"`)
	assert.NotContains(t, rec.Body.String(), `"documentationUrl": "https://compose.com/databases/scylladb"`)
	assert.NotContains(t, rec.Body.String(), `"imageUrl": "https://compose.com/assets/icd-icons/etcd-9bf4cedacb096e58868085dc0b91b8c3cc3c1f0b3f8be05d1bf7ea5b5e9b6697.svg"`)
	assert.Contains(t, rec.Body.String(), `"2 GB Storage"`)
	assert.NotContains(t, rec.Body.String(), `"5 GB Storage"`)
	assert.Contains(t, rec.Body.String(), `"longDescription": "Deploy RabbitMQ on AWS, GCP, or IBM Cloud in minutes. Fully managed, highly-available and production ready."`)
	assert.Equal(t, util.Body("../_fixtures/broker_catalog_trimmed.json"), rec.Body.String())
}
