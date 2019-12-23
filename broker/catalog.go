package broker

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/JamesClonk/compose-broker/log"
	yaml "gopkg.in/yaml.v2"
)

type ServiceCatalog struct {
	Services []Service `json:"services" yaml:"services"`
}
type Service struct {
	ID                   string   `json:"id" yaml:"id"`
	Name                 string   `json:"name" yaml:"name"`
	Description          string   `json:"description" yaml:"description"`
	Bindable             bool     `json:"bindable" yaml:"bindable"`
	InstancesRetrievable bool     `json:"instances_retrievable" yaml:"instances_retrievable"`
	BindingsRetrievable  bool     `json:"bindings_retrievable" yaml:"bindings_retrievable"`
	PlanUpdateable       bool     `json:"plan_updateable" yaml:"plan_updateable"`
	Tags                 []string `json:"tags" yaml:"tags"`
	Metadata             struct {
		DisplayName         string `json:"displayName" yaml:"displayName"`
		ImageURL            string `json:"imageUrl,omitempty" yaml:"imageUrl,omitempty"`
		LongDescription     string `json:"longDescription" yaml:"longDescription"`
		ProviderDisplayName string `json:"providerDisplayName" yaml:"providerDisplayName"`
		DocumentationURL    string `json:"documentationUrl" yaml:"documentationUrl"`
		SupportURL          string `json:"supportUrl" yaml:"supportUrl"`
	} `json:"metadata" yaml:"metadata"`
	Plans []ServicePlan `json:"plans" yaml:"plans"`
}
type ServicePlan struct {
	ID          string `json:"id" yaml:"id"`
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
	Free        bool   `json:"free" yaml:"free"`
	Bindable    bool   `json:"bindable" yaml:"bindable"`
	Metadata    struct {
		DisplayName string `json:"displayName" yaml:"displayName"`
		ImageURL    string `json:"imageUrl,omitempty" yaml:"imageUrl,omitempty"`
		Costs       []struct {
			Amount struct {
				USDollar float64 `json:"usd" yaml:"usd"`
			} `json:"amount" yaml:"amount"`
			Unit string `json:"unit" yaml:"unit"`
		} `json:"costs" yaml:"costs"`
		Bullets          []string `json:"bullets" yaml:"bullets"`
		HighAvailability bool     `json:"highAvailability" yaml:"highAvailability"`
		Units            int      `json:"units" yaml:"units"`
		CacheMode        bool     `json:"cache_mode,omitempty" yaml:"cache_mode,omitempty"`
		Version          string   `json:"version,omitempty" yaml:"version,omitempty"`
	} `json:"metadata" yaml:"metadata"`
}

func LoadServiceCatalog(filename string) *ServiceCatalog {
	var catalog ServiceCatalog

	if _, err := os.Stat(filename); err == nil {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Errorf("could not load %s", filename)
			log.Fatalln(err)
		}
		if err := yaml.Unmarshal(data, &catalog); err != nil {
			log.Errorf("could not parse %s", filename)
			log.Fatalln(err)
		}

		// expect & hardcode certain default values
		if len(catalog.Services) < 1 {
			log.Errorln("invalid catalog, no service offerings defined")
			log.Fatalln(catalog)
		}
		for sx, service := range catalog.Services {
			if len(service.ID) == 0 {
				log.Errorf("service #%d: ID is missing in catalog %s", sx, filename)
				log.Fatalln(catalog)
			}
			if len(service.Name) == 0 {
				log.Errorf("service #%d: name is missing in catalog %s", sx, filename)
				log.Fatalln(catalog)
			}
			// displayName
			if len(service.Metadata.DisplayName) == 0 {
				catalog.Services[sx].Metadata.DisplayName = service.Name
			}
			// enforce flags
			catalog.Services[sx].Bindable = true
			catalog.Services[sx].InstancesRetrievable = true
			catalog.Services[sx].BindingsRetrievable = true
			// catalog.Services[sx].PlanUpdateable = true // don't enforce "plan_updateable", some databases might truly not support scaling

			if len(service.Plans) < 1 {
				log.Errorln("invalid catalog, at least one service plan has to be defined")
				log.Fatalln(catalog)
			}
			for px, plan := range service.Plans {
				if len(plan.ID) == 0 {
					log.Errorf("service #%d, plan #%d: ID is missing in catalog %s", sx, px, filename)
					log.Fatalln(catalog)
				}
				if len(plan.Name) == 0 {
					log.Errorf("service #%d, plan #%d: name is missing in catalog %s", sx, px, filename)
					log.Fatalln(catalog)
				}
				if plan.Metadata.Units == 0 {
					catalog.Services[sx].Plans[px].Metadata.Units = 1
				}
			}
		}

	} else {
		log.Errorf("could not load %s", filename)
		log.Fatalln(err)
	}
	return &catalog
}

func (b *Broker) Catalog(rw http.ResponseWriter, req *http.Request) {
	// filter catalog by /databases api response, trim everything that is not at least "stable" or "beta"
	databases, err := b.Client.GetDatabases()
	if err != nil {
		log.Errorf("could not filter services for catalog: %v", err)
	}

	filteredServices := make([]Service, 0)
	for _, service := range b.ServiceCatalog.Services {
		for _, database := range databases {
			if service.Name == database.DatabaseType {
				// only allow stable or beta service offerings
				if database.Status == "stable" || database.Status == "beta" {
					filteredServices = append(filteredServices, service)
				}
			}
		}
	}
	b.ServiceCatalog.Services = filteredServices

	b.write(rw, req, 200, b.ServiceCatalog)
}
