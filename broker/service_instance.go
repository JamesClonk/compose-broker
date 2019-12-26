package broker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/JamesClonk/compose-broker/api"
	"github.com/JamesClonk/compose-broker/log"
	"github.com/gorilla/mux"
)

type ServiceInstanceProvisioning struct {
	ServiceID  string `json:"service_id"`
	PlanID     string `json:"plan_id"`
	Parameters struct {
		AccountID  string `json:"account_id"`
		Datacenter string `json:"datacenter"`
		Version    string `json:"version"`
		Units      int    `json:"units"`
		CacheMode  bool   `json:"cache_mode"`
	} `json:"parameters"`
}
type ServiceInstanceProvisioningResponse struct {
	DashboardURL string `json:"dashboard_url"`
}

type ServiceInstanceFetchResponse struct {
	DashboardURL string                                 `json:"dashboard_url"`
	Parameters   ServiceInstanceFetchResponseParameters `json:"parameters"`
}
type ServiceInstanceFetchResponseParameters struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	AccountID      string    `json:"account_id"`
	Type           string    `json:"type"`
	Notes          string    `json:"notes"`
	Version        string    `json:"version"`
	CreatedAt      time.Time `json:"created_at"`
	AllocatedUnits int       `json:"allocated_units"`
	UsedUnits      int       `json:"used_units"`
}

type ServiceInstanceUpdate struct {
	ServiceID  string `json:"service_id"`
	PlanID     string `json:"plan_id"`
	Parameters struct {
		Units int `json:"units"`
	} `json:"parameters"`
}
type ServiceInstanceUpdateResponse struct {
	DashboardURL string `json:"dashboard_url"`
}

func (b *Broker) ProvisionInstance(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	instanceID := vars["instanceID"]

	// verify request is async, must have query param "?accepts_incomplete=true"
	incomplete := req.URL.Query().Get("accepts_incomplete")
	if incomplete != "true" {
		log.Errorf("creating service instance %s requires async / accepts_incomplete=true", instanceID)
		b.Error(rw, req, 422, "AsyncRequired", "Service instance provisioning requires an asynchronous operation")
		return
	}

	if req.Body == nil {
		log.Errorf("error reading provisioning request for service instance %s: %v", instanceID, req)
		b.Error(rw, req, 400, "MalformedRequest", "Could not read provisioning request")
		return
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Errorln(err)
		log.Errorf("error reading provisioning request for service instance %s: %v", instanceID, req)
		b.Error(rw, req, 400, "MalformedRequest", "Could not read provisioning request")
		return
	}
	if len(body) == 0 {
		body = []byte("{}")
	}

	var provisioning ServiceInstanceProvisioning
	if err := json.Unmarshal([]byte(body), &provisioning); err != nil {
		log.Errorln(err)
		log.Errorf("could not unmarshal provisioning request body for service instance %s: %v", instanceID, string(body))
		b.Error(rw, req, 400, "MalformedRequest", "Could not unmarshal provisioning request")
		return
	}

	// collect deployment values
	var deploymentType, accountID, datacenter, version string
	var cacheMode bool
	var units int
	if len(provisioning.PlanID) > 0 {
		// get plan values
		for _, service := range b.ServiceCatalog.Services {
			if provisioning.ServiceID == service.ID {
				for _, plan := range service.Plans {
					if provisioning.PlanID == plan.ID {
						deploymentType = service.Name
						datacenter = plan.Metadata.Datacenter
						version = plan.Metadata.Version
						units = plan.Metadata.Units
						cacheMode = plan.Metadata.CacheMode
					}
				}
			}
		}
		if len(deploymentType) == 0 {
			log.Errorf("could not find plan_id %s for provisioning service instance %s", provisioning.PlanID, instanceID)
			b.Error(rw, req, 400, "MalformedRequest", "Unknown plan_id")
			return
		}
	}

	// units can also be provided as provisioning parameter, takes precedence over plan value
	if provisioning.Parameters.Units > 0 {
		units = provisioning.Parameters.Units
	}
	// verify scaling target value (units)
	if units < 1 {
		log.Errorf("units value %d must be greater than 0 for provisioning service instance %s", units, instanceID)
		b.Error(rw, req, 400, "MissingParameters", "Units parameter is missing for service instance provisioning")
		return
	}

	// account_id can be set to a global default value
	if len(b.APIConfig.DefaultAccountID) > 0 {
		accountID = b.APIConfig.DefaultAccountID
	}
	// account_id can be provided as provisioning parameter
	if len(provisioning.Parameters.AccountID) > 0 {
		accountID = provisioning.Parameters.AccountID
	}
	if len(accountID) == 0 {
		// get accountID from API
		accounts, err := b.Client.GetAccounts()
		if err != nil {
			log.Errorf("could not fetch accounts: %v", err)
			b.Error(rw, req, 409, "UnknownError", "Could not read Compose.io accounts") // TODO: write test case
			return
		}
		if len(accounts) > 0 {
			accountID = accounts[0].ID // TODO: write test case
		}
	}
	if len(accountID) == 0 {
		log.Errorf("account_id for provisioning service instance %s could not be determined", instanceID)
		b.Error(rw, req, 400, "MissingParameters", "AccountID is missing for service instance provisioning")
		return
	}

	// datacenter can also be provided as provisioning parameter, takes precedence over plan value
	if len(provisioning.Parameters.Datacenter) > 0 {
		datacenter = provisioning.Parameters.Datacenter
	}
	if len(datacenter) == 0 {
		// get default datacenter as fallback
		datacenter = b.APIConfig.DefaultDatacenter
	}

	// version can also be provided as provisioning parameter, takes precedence over plan value
	if len(provisioning.Parameters.Version) > 0 {
		version = provisioning.Parameters.Version
	}

	// cache_mode can also be provided as provisioning parameter, takes precedence over plan value
	if provisioning.Parameters.CacheMode {
		cacheMode = provisioning.Parameters.CacheMode
	}

	// check if it already exists
	instance, err := b.Client.GetDeploymentByName(instanceID)
	if err == nil && instance.Name == instanceID {
		recipes, err := b.Client.GetRecipes(instance.ID)
		if err != nil {
			log.Warnf("could not fetch any recipes for service instance %s: %v", instanceID, err)
		}
		if len(recipes) > 0 {
			recipes.SortByUpdatedAt()

			// response JSON
			provisionResponse := ServiceInstanceProvisioningResponse{
				DashboardURL: strings.TrimSuffix(instance.Links.ComposeWebUI.HREF, "{?embed}"),
			}
			if recipes[0].Name == "Provision" &&
				(recipes[0].Status == "running" ||
					recipes[0].Status == "waiting") {
				log.Infof("service instance %s is already ongoing provisioning, nothing to do", instanceID)
				b.write(rw, req, 202, provisionResponse) // TODO: write test case
				return
			}
			if recipes[0].Status == "complete" {
				scaling, err := b.Client.GetScaling(instance.ID)
				if err == nil && scaling.AllocatedUnits == units {
					log.Infof("service instance %s already exists and has same scaling, nothing to do", instanceID)
					b.write(rw, req, 200, provisionResponse) // TODO: write test case
					return
				}
			}
		}
		log.Errorf("could not create service instance %s: %v", instanceID, err)
		b.Error(rw, req, 409, "UnknownError", "Could not create service instance") // TODO: write test case
		return
	}

	// provision service instance
	newDeployment := api.NewDeployment{
		Name:       instanceID,
		AccountID:  accountID,
		Datacenter: datacenter,
		Type:       deploymentType,
		Version:    version,
		Units:      units,
		CacheMode:  cacheMode,
		Notes:      fmt.Sprintf("%s-%s", provisioning.ServiceID, provisioning.PlanID),
	}
	deployment, err := b.Client.CreateDeployment(newDeployment)
	if err != nil {
		log.Errorf("could not create service instance %s: %v", instanceID, err)
		b.Error(rw, req, 500, "UnknownError", "Could not provision service instance") // TODO: write test case
		return
	}

	if len(deployment.ProvisionRecipeID) > 0 {
		if state, err := b.Client.GetRecipe(deployment.ProvisionRecipeID); err == nil {
			if state.Status == "complete" {
				// TODO: write test case
				b.write(rw, req, 201, map[string]string{}) // provisioning already done
				return
			} else if state.Status == "failed" {
				// TODO: write test case
				log.Errorf("could not create service instance %s, recipe %s failed", instanceID, deployment.ProvisionRecipeID)
				b.Error(rw, req, 400, "ProvisionFailure", "Could not create service instance") // provisioning immediately failed
				return
			}
		}
	}

	// response JSON
	provisionResponse := ServiceInstanceProvisioningResponse{
		DashboardURL: strings.TrimSuffix(deployment.Links.ComposeWebUI.HREF, "{?embed}"),
	}
	b.write(rw, req, 202, provisionResponse) // default async response
}

func (b *Broker) FetchInstance(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	instanceID := vars["instanceID"]

	instance, err := b.Client.GetDeploymentByName(instanceID)
	if err != nil || instance.Name != instanceID {
		log.Errorf("could not fetch service instance %s: %v", instanceID, err)
		b.Error(rw, req, 404, "MissingServiceInstance", "The service instance does not exist")
		return
	}

	recipes, err := b.Client.GetRecipes(instance.ID)
	if err != nil {
		log.Errorf("could not fetch recipes for service instance %s: %v", instanceID, err)
		b.Error(rw, req, 404, "MissingRecipes", "The service instance recipes could not be found")
		return
	}
	if len(recipes) > 0 {
		recipes.SortByUpdatedAt()
		if recipes[0].Status == "running" ||
			recipes[0].Status == "waiting" {
			log.Warnf("service instance %s has currently an ongoing recipe", instanceID)
			if recipes[0].Name == "Provision" {
				b.Error(rw, req, 404, "ConcurrencyError", "The service instance provisioning is still in progress")
				return
			} else {
				b.Error(rw, req, 422, "ConcurrencyError", "The service instance is being updated")
				return
			}
		}
	}

	scaling, err := b.Client.GetScaling(instance.ID)
	if err != nil {
		log.Errorf("could not fetch scaling parameters for service instance %s: %v", instanceID, err)
		b.Error(rw, req, 404, "MissingScalingParameters", "The service instance scaling parameters do not exist")
		return
	}

	// response JSON
	fetchResponse := ServiceInstanceFetchResponse{
		DashboardURL: strings.TrimSuffix(instance.Links.ComposeWebUI.HREF, "{?embed}"),
		Parameters: ServiceInstanceFetchResponseParameters{
			ID:             instance.ID,
			Name:           instance.Name,
			AccountID:      instance.AccountID,
			Type:           instance.Type,
			Notes:          instance.Notes,
			Version:        instance.Version,
			CreatedAt:      instance.CreatedAt,
			AllocatedUnits: scaling.AllocatedUnits,
			UsedUnits:      scaling.UsedUnits,
		},
	}
	b.write(rw, req, 200, fetchResponse)
}

func (b *Broker) UpdateInstance(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	instanceID := vars["instanceID"]

	// verify request is async, must have query param "?accepts_incomplete=true"
	incomplete := req.URL.Query().Get("accepts_incomplete")
	if incomplete != "true" {
		log.Errorf("updating service instance %s requires async / accepts_incomplete=true", instanceID)
		b.Error(rw, req, 422, "AsyncRequired", "Service instance updating requires an asynchronous operation")
		return
	}

	if req.Body == nil {
		log.Errorf("error reading update request for service instance %s: %v", instanceID, req)
		b.Error(rw, req, 400, "MalformedRequest", "Could not read update request")
		return
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Errorln(err)
		log.Errorf("error reading update request for service instance %s: %v", instanceID, req)
		b.Error(rw, req, 400, "MalformedRequest", "Could not read update request")
		return
	}
	if len(body) == 0 {
		body = []byte("{}")
	}

	var update ServiceInstanceUpdate
	if err := json.Unmarshal([]byte(body), &update); err != nil {
		log.Errorln(err)
		log.Errorf("could not unmarshal update request body for service instance %s: %v", instanceID, string(body))
		b.Error(rw, req, 400, "MalformedRequest", "Could not unmarshal update request")
		return
	}

	// verify scaling target value (units), by either taking value of plan or by provided parameter
	var units int
	if len(update.PlanID) > 0 {
		// get units if plan was specified
		for _, service := range b.ServiceCatalog.Services {
			if update.ServiceID == service.ID {
				for _, plan := range service.Plans {
					if update.PlanID == plan.ID {
						units = plan.Metadata.Units
					}
				}
			}
		}
		if units == 0 {
			log.Errorf("could not find plan_id %s for updating service instance %s", update.PlanID, instanceID)
			b.Error(rw, req, 400, "MalformedRequest", "Unknown plan_id")
			return
		}
	}
	if update.Parameters.Units > 0 {
		units = update.Parameters.Units
	}
	if units < 1 {
		log.Errorf("units value %d must be greater than 0 for updating service instance %s", units, instanceID)
		b.Error(rw, req, 400, "MissingParameters", "Units parameter is missing for service instance update")
		return
	}

	instance, err := b.Client.GetDeploymentByName(instanceID)
	if err != nil || instance.Name != instanceID {
		log.Errorf("could not fetch service instance %s: %v", instanceID, err)
		b.Error(rw, req, 404, "ServiceInstanceNotFound", "The service instance does not exist")
		return
	}

	// would it actually do anything?
	scaling, err := b.Client.GetScaling(instance.ID)
	if err != nil {
		log.Errorf("could not fetch scaling parameters for service instance %s: %v", instanceID, err)
		b.Error(rw, req, 409, "UnknownError", "Could not read service instance scaling")
		return
	}
	if scaling.AllocatedUnits == units {
		log.Warnf("service instance %s already has %d units", instanceID, units)
		b.write(rw, req, 200, map[string]string{}) // update would have no effect
		return
	}

	// return concurrency error if there is still/already another recipe ongoing for this deployment
	recipes, err := b.Client.GetRecipes(instance.ID)
	if err != nil {
		log.Warnf("could not fetch any recipes for service instance %s: %v", instanceID, err)
	}
	if len(recipes) > 0 {
		recipes.SortByUpdatedAt()
		if recipes[0].Status == "running" ||
			recipes[0].Status == "waiting" {
			log.Errorf("updating service instance %s not possible due to an ongoing recipe", instanceID)
			b.Error(rw, req, 422, "ConcurrencyError", "The service instance is currently being updated")
			return
		}
	}

	recipe, err := b.Client.UpdateScaling(instance.ID, units)
	if err != nil {
		log.Errorf("could not update service instance %s: %v", instanceID, err)
		b.Error(rw, req, 409, "UnknownError", "Could not update service instance")
		return
	}

	if len(recipe.ID) > 0 {
		if state, err := b.Client.GetRecipe(recipe.ID); err == nil {
			if state.Status == "complete" {
				b.write(rw, req, 200, map[string]string{}) // update already done
				return
			} else if state.Status == "failed" {
				log.Errorf("could not update service instance %s, recipe %s failed", instanceID, recipe.ID)
				b.Error(rw, req, 409, "UpdateFailure", "Could not update service instance") // update immediately failed
				return
			}
		}
	}

	// response JSON
	updateResponse := ServiceInstanceUpdateResponse{
		DashboardURL: strings.TrimSuffix(instance.Links.ComposeWebUI.HREF, "{?embed}"),
	}
	b.write(rw, req, 202, updateResponse) // default async response
}

func (b *Broker) DeprovisionInstance(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	instanceID := vars["instanceID"]

	// verify request is async, must have query param "?accepts_incomplete=true"
	incomplete := req.URL.Query().Get("accepts_incomplete")
	if incomplete != "true" {
		log.Errorf("deleting service instance %s requires async / accepts_incomplete=true", instanceID)
		b.Error(rw, req, 422, "AsyncRequired", "Service instance deprovisioning requires an asynchronous operation")
		return
	}

	instance, err := b.Client.GetDeploymentByName(instanceID)
	if err != nil || instance.Name != instanceID {
		log.Errorf("could not find service instance %s: %v", instanceID, err)
		b.Error(rw, req, 410, "MissingServiceInstance", "The service instance does not exist")
		return
	}

	// return concurrency error if there is still/already another recipe ongoing for this deployment
	recipes, err := b.Client.GetRecipes(instance.ID)
	if err != nil {
		log.Warnf("could not fetch any recipes for service instance %s: %v", instanceID, err)
	}
	if len(recipes) > 0 {
		recipes.SortByUpdatedAt()
		if recipes[0].Status == "running" ||
			recipes[0].Status == "waiting" {
			log.Errorf("deleting service instance %s not possible due to an ongoing recipe", instanceID)
			b.Error(rw, req, 422, "ConcurrencyError", "The service instance is currently being updated")
			return
		}
	}

	// deprovision service instance
	recipe, err := b.Client.DeleteDeployment(instance.ID)
	if err != nil {
		log.Errorf("could not delete service instance %s: %v", instanceID, err)
		b.Error(rw, req, 500, "UnknownError", "Could not delete service instance")
		return
	}

	if len(recipe.ID) > 0 {
		if state, err := b.Client.GetRecipe(recipe.ID); err == nil {
			if state.Status == "complete" {
				b.write(rw, req, 200, map[string]string{}) // deletion already done
				return
			} else if state.Status == "failed" {
				log.Errorf("could not delete service instance %s, recipe %s failed", instanceID, recipe.ID)
				b.Error(rw, req, 500, "DeprovisionFailure", "Could not delete service instance") // deletion immediately failed
				return
			}
		}
	}
	b.write(rw, req, 202, map[string]string{}) // default async response
}
