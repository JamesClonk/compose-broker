package broker

import (
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/JamesClonk/compose-broker/api"
	"github.com/JamesClonk/compose-broker/log"
	"github.com/gorilla/mux"
)

type ServiceBindingResponse struct {
	Credentials ServiceBindingResponseCredentials `json:"credentials"`
	Endpoints   []ServiceBindingResponseEndpoint  `json:"endpoints"`
	Parameters  ServiceBindingResponseParameters  `json:"parameters"`
}
type ServiceBindingResponseCredentials struct {
	Direct        []string `json:"direct"`
	CLI           []string `json:"cli"`
	Maps          []string `json:"maps"`
	SSH           []string `json:"ssh"`
	Health        []string `json:"health"`
	Admin         []string `json:"admin"`
	URI           string   `json:"uri,omitempty"`
	URL           string   `json:"url,omitempty"`
	DatabaseURI   string   `json:"database_uri,omitempty"`
	Username      string   `json:"username,omitempty"`
	Password      string   `json:"password,omitempty"`
	Database      string   `json:"database,omitempty"`
	Scheme        string   `json:"scheme,omitempty"`
	Host          string   `json:"host,omitempty"`
	Hostname      string   `json:"hostname,omitempty"`
	Port          int      `json:"port,omitempty"`
	CACertificate string   `json:"ca_certificate,omitempty"`
}
type ServiceBindingResponseEndpoint struct {
	Host  string   `json:"host"`
	Ports []string `json:"ports"`
}
type ServiceBindingResponseParameters struct {
	Deployment api.Deployment `json:"deployment"`
	Scaling    api.Scaling    `json:"scaling,omitempty"`
}

func (b *Broker) FetchBinding(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	instanceID := vars["instanceID"]

	instance, err := b.Client.GetDeploymentByName(instanceID)
	if err != nil || instance.Name != instanceID {
		log.Errorf("could not query service instance %s: %v", instanceID, err)
		b.Error(rw, req, 404, "MissingServiceInstance", "The service instance does not exist")
		return
	}
	b.write(rw, req, 200, b.getBinding(instance))
}

func (b *Broker) Unbind(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	instanceID := vars["instanceID"]

	instance, err := b.Client.GetDeploymentByName(instanceID)
	if err != nil || instance.Name != instanceID {
		log.Errorf("could not query service instance %s: %v", instanceID, err)
		b.Error(rw, req, 410, "MissingServiceInstance", "The service instance does not exist")
		return
	}
	b.write(rw, req, 200, map[string]string{})
}

func (b *Broker) getBinding(deployment *api.Deployment) ServiceBindingResponse {
	credentials := ServiceBindingResponseCredentials{
		Direct:        deployment.ConnectionStrings.Direct,
		CLI:           deployment.ConnectionStrings.CLI,
		Maps:          deployment.ConnectionStrings.Maps,
		SSH:           deployment.ConnectionStrings.SSH,
		Health:        deployment.ConnectionStrings.Health,
		Admin:         deployment.ConnectionStrings.Admin,
		CACertificate: deployment.CACertificateBase64,
	}
	endpoints := make([]ServiceBindingResponseEndpoint, 0)

	if len(deployment.ConnectionStrings.Direct) > 0 {
		credentials.URI = deployment.ConnectionStrings.Direct[0]
		credentials.URL = deployment.ConnectionStrings.Direct[0]
		credentials.DatabaseURI = deployment.ConnectionStrings.Direct[0]
	}

	// get credentials
	if u, err := url.Parse(credentials.URI); err == nil {
		credentials.Username = u.User.Username()
		password, _ := u.User.Password()
		credentials.Password = password

		credentials.Scheme = u.Scheme
		credentials.Host = u.Host
		hostname, port, _ := net.SplitHostPort(u.Host)
		credentials.Hostname = hostname
		p, _ := strconv.Atoi(port)
		credentials.Port = p

		credentials.Database = strings.TrimPrefix(u.Path, "/")
		if strings.Contains(credentials.Database, "?") {
			rx := regexp.MustCompile(`([^\?]*)\?.*`) // trim connection options
			credentials.Database = rx.ReplaceAllString(credentials.Database, "${1}")
		}

	}

	// get endpoints
	for _, connstr := range deployment.ConnectionStrings.Direct {
		if u, err := url.Parse(connstr); err == nil {
			hostname, port, _ := net.SplitHostPort(u.Host)

			endpoints = append(endpoints, ServiceBindingResponseEndpoint{
				Host:  hostname,
				Ports: []string{port},
			})
		}
	}

	scaling, err := b.Client.GetScaling(deployment.ID)
	if err != nil {
		log.Warnf("could not query scaling parameters for service instance %s: %v", deployment.ID, err)
		scaling = &api.Scaling{}
	}

	return ServiceBindingResponse{
		Credentials: credentials,
		Endpoints:   endpoints,
		Parameters: ServiceBindingResponseParameters{
			Deployment: *deployment,
			Scaling:    *scaling,
		},
	}
}
