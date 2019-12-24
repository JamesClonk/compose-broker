package broker

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/JamesClonk/compose-broker/log"
	"github.com/gorilla/mux"
)

type ServiceInstanceProvisioning struct {
	ServiceID  string `json:"service_id"`
	PlanID     string `json:"plan_id"`
	Parameters struct {
		Units int `json:"units"`
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

func (b *Broker) FetchInstance(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	instanceID := vars["instanceID"]

	instance, err := b.Client.GetDeploymentByName(instanceID)
	if err != nil || instance.Name != instanceID {
		log.Errorf("could not fetch service instance: %v", err)
		b.Error(rw, req, 404, "MissingServiceInstance", "The service instance does not exist")
		return
	}

	recipes, err := b.Client.GetRecipes(instance.ID)
	if err != nil {
		log.Errorf("could not fetch service instance recipes: %v", err)
		b.Error(rw, req, 404, "MissingRecipes", "The service instance recipes could not be found")
		return
	}
	if len(recipes) > 0 {
		recipes.SortByUpdatedAt()
		if recipes[0].Status == "running" ||
			recipes[0].Status == "waiting" {
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
		log.Errorf("could not fetch service instance scaling parameters: %v", err)
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
	// verify request is async, must have query param "?accepts_incomplete=true"
	incomplete := req.URL.Query().Get("accepts_incomplete")
	if incomplete != "true" {
		b.Error(rw, req, 422, "AsyncRequired", "Service instance updating requires an asynchronous operation")
		return
	}

	if req.Body == nil {
		log.Errorf("error reading update request: %v", req)
		b.Error(rw, req, 400, "MalformedRequest", "Could not read update request")
		return
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Errorln(err)
		log.Errorf("error reading update request: %v", req)
		b.Error(rw, req, 400, "MalformedRequest", "Could not read update request")
		return
	}
	if len(body) == 0 {
		body = []byte("{}")
	}

	var update ServiceInstanceUpdate
	if err := json.Unmarshal([]byte(body), &update); err != nil {
		log.Errorln(err)
		log.Errorf("could not unmarshal update request body: %v", string(body))
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
			log.Errorf("could not find plan_id [%s] for service instance update", update.PlanID)
			b.Error(rw, req, 400, "MalformedRequest", "Unknown plan_id")
			return
		}
	}
	if update.Parameters.Units > 0 {
		units = update.Parameters.Units
	}
	if units < 1 {
		log.Errorf("units value %d must be greater than 0 for service instance update", units)
		b.Error(rw, req, 400, "MissingParameters", "Units parameter is missing for service instance update")
		return
	}

	vars := mux.Vars(req)
	instanceID := vars["instanceID"]

	instance, err := b.Client.GetDeploymentByName(instanceID)
	if err != nil || instance.Name != instanceID {
		log.Errorf("could not fetch service instance: %v", err)
		b.Error(rw, req, 404, "ServiceInstanceNotFound", "The service instance does not exist")
		return
	}

	// would it actually do anything?
	scaling, err := b.Client.GetScaling(instance.ID)
	if err != nil {
		log.Errorf("could not fetch service instance scaling parameters: %v", err)
		b.Error(rw, req, 409, "UnknownError", "Could not read service instance scaling")
		return
	}
	if scaling.AllocatedUnits == units {
		log.Warnf("service instance already has %d units", units)
		b.write(rw, req, 200, map[string]string{}) // update would have no effect
		return
	}

	// return concurrency error if there is still/already another recipe ongoing for this deployment
	recipes, err := b.Client.GetRecipes(instance.ID)
	if err != nil {
		log.Warnf("could not fetch any service instance recipes: %v", err)
	}
	if len(recipes) > 0 {
		recipes.SortByUpdatedAt()
		if recipes[0].Status == "running" ||
			recipes[0].Status == "waiting" {
			b.Error(rw, req, 422, "ConcurrencyError", "The service instance is currently being updated")
			return
		}
	}

	recipe, err := b.Client.UpdateScaling(instance.ID, units)
	if err != nil {
		log.Errorf("could not update service instance: %v", err)
		b.Error(rw, req, 409, "UnknownError", "Could not update service instance")
		return
	}

	if len(recipe.ID) > 0 {
		if state, err := b.Client.GetRecipe(recipe.ID); err == nil {
			if state.Status == "complete" {
				b.write(rw, req, 200, map[string]string{}) // update already done
				return
			} else if state.Status == "failed" {
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
	// verify request is async, must have query param "?accepts_incomplete=true"
	incomplete := req.URL.Query().Get("accepts_incomplete")
	if incomplete != "true" {
		b.Error(rw, req, 422, "AsyncRequired", "Service instance deprovisioning requires an asynchronous operation")
		return
	}

	vars := mux.Vars(req)
	instanceID := vars["instanceID"]

	instance, err := b.Client.GetDeploymentByName(instanceID)
	if err != nil || instance.Name != instanceID {
		log.Errorf("could not find service instance: %v", err)
		b.Error(rw, req, 410, "MissingServiceInstance", "The service instance does not exist")
		return
	}

	// return concurrency error if there is still/already another recipe ongoing for this deployment
	recipes, err := b.Client.GetRecipes(instance.ID)
	if err != nil {
		log.Warnf("could not fetch any service instance recipes: %v", err)
	}
	if len(recipes) > 0 {
		recipes.SortByUpdatedAt()
		if recipes[0].Status == "running" ||
			recipes[0].Status == "waiting" {
			b.Error(rw, req, 422, "ConcurrencyError", "The service instance is currently being updated")
			return
		}
	}

	// deprovision service instance
	recipe, err := b.Client.DeleteDeployment(instance.ID)
	if err != nil {
		log.Errorf("could not delete service instance: %v", err)
		b.Error(rw, req, 500, "UnknownError", "Could not delete service instance")
		return
	}

	if len(recipe.ID) > 0 {
		if state, err := b.Client.GetRecipe(recipe.ID); err == nil {
			if state.Status == "complete" {
				b.write(rw, req, 200, map[string]string{}) // deletion already done
				return
			} else if state.Status == "failed" {
				b.Error(rw, req, 500, "DeprovisionFailure", "Could not delete service instance") // deletion immediately failed
				return
			}
		}
	}
	b.write(rw, req, 202, map[string]string{}) // default async response
}
