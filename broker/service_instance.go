package broker

import (
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
		Region string `json:"region"`
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

func (b *Broker) FetchInstance(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	instanceID := vars["instanceID"]

	instance, err := b.Client.GetDeploymentByName(instanceID)
	if err != nil || instance.Name != instanceID {
		log.Errorf("could not fetch service instance: %v", err)
		b.Error(rw, req, 404, "ServiceInstanceNotFound", "The service instance does not exist")
		return
	}

	recipes, err := b.Client.GetRecipes(instance.ID)
	if err != nil {
		log.Errorf("could not fetch service instance recipes: %v", err)
		b.Error(rw, req, 404, "RecipesNotFound", "The service instance recipes could not be found")
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
		b.Error(rw, req, 404, "ScalingParametersNotFound", "The service instance scaling parameters do not exist")
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
