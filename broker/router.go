package broker

import (
	"github.com/JamesClonk/compose-broker/config"
	"github.com/gorilla/mux"
)

func NewRouter(c *config.Config) *mux.Router {
	b := NewBroker(c)

	// mux router
	r := mux.NewRouter()

	// routes
	r.HandleFunc("/", b.BasicAuth(b.Health)).Methods("GET")
	r.PathPrefix("/health").HandlerFunc(b.Health)

	r.HandleFunc("/v2/catalog", b.BasicAuth(b.Catalog)).Methods("GET")

	r.HandleFunc("/v2/service_instances/{instanceID}", b.BasicAuth(b.ProvisionInstance)).Methods("PUT")
	r.HandleFunc("/v2/service_instances/{instanceID}/last_operation", b.BasicAuth(b.LastOperationOnInstance)).Methods("GET")
	r.HandleFunc("/v2/service_instances/{instanceID}", b.BasicAuth(b.FetchInstance)).Methods("GET")
	r.HandleFunc("/v2/service_instances/{instanceID}", b.BasicAuth(b.UpdateInstance)).Methods("PATCH")
	r.HandleFunc("/v2/service_instances/{instanceID}", b.BasicAuth(b.DeprovisionInstance)).Methods("DELETE")

	// TODO: bindings...
	// r.HandleFunc("/v2/service_instances/{instanceID}/service_bindings/{bindingID}", b.BasicAuth(b.Bind)).Methods("PUT")
	r.HandleFunc("/v2/service_instances/{instanceID}/service_bindings/{bindingID}", b.BasicAuth(b.FetchBinding)).Methods("GET")
	r.HandleFunc("/v2/service_instances/{instanceID}/service_bindings/{bindingID}", b.BasicAuth(b.Unbind)).Methods("DELETE")

	return r
}
