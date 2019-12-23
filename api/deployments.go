package api

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/JamesClonk/compose-broker/log"
)

type Deployments []Deployment
type Deployment struct {
	ID                  string    `json:"id"`
	AccountID           string    `json:"account_id"`
	Name                string    `json:"name"`
	Type                string    `json:"type"`
	Notes               string    `json:"notes"`
	CustomerBillingCode string    `json:"customer_billing_code"`
	ClusterID           string    `json:"cluster_id"`
	Version             string    `json:"version"`
	CACertificateBase64 string    `json:"ca_certificate_base64"`
	ProvisionRecipeID   string    `json:"provision_recipe_id"`
	CreatedAt           time.Time `json:"created_at"`
	ConnectionStrings   struct {
		Direct []string `json:"direct"`
		CLI    []string `json:"cli"`
		Maps   []string `json:"maps"`
		SSH    []string `json:"ssh"`
		Health []string `json:"health"`
		Admin  []string `json:"admin"`
	} `json:"connection_strings"`
	Links struct {
		ComposeWebUI struct {
			HREF      string `json:"href"`
			Templated bool   `json:"templated"`
		} `json:"compose_web_ui"`
	} `json:"_links"`
}
type NewDeployment struct {
	Name       string `json:"name"`
	AccountID  string `json:"account_id"`
	Datacenter string `json:"datacenter"`
	Type       string `json:"type"`
	Version    string `json:"version,omitempty"`
	Units      int    `json:"units,omitempty"`
	CacheMode  bool   `json:"cache_mode,omitempty"`
	Notes      string `json:"notes,omitempty"`
}

func (c *Client) CreateDeployment(newDeployment NewDeployment) (*Deployment, error) {
	// set defaults
	if len(newDeployment.Datacenter) == 0 {
		newDeployment.Datacenter = c.Config.DefaultDatacenter
	}
	if newDeployment.Units < 1 {
		newDeployment.Units = 1
	}

	data := struct {
		Deployment NewDeployment `json:"deployment"`
	}{
		Deployment: newDeployment,
	}
	payload, err := json.Marshal(data)
	if err != nil {
		log.Errorf("could not marshal deployment payload: %#v", newDeployment)
		return nil, err
	}

	body, err := c.PostAsync("deployments", string(payload))
	if err != nil {
		log.Errorf("could not create Compose.io deployment: %s", err)
		return nil, err
	}

	deployment := &Deployment{}
	if err := json.Unmarshal([]byte(body), deployment); err != nil {
		log.Errorf("could not unmarshal deployment response: %#v", body)
		return nil, err
	}
	return deployment, nil
}

func (c *Client) GetDeployments() (Deployments, error) {
	body, err := c.Get("deployments")
	if err != nil {
		log.Errorf("could not get Compose.io deployments: %s", err)
		return nil, err
	}

	response := struct {
		Embedded struct {
			Deployments Deployments `json:"deployments"`
		} `json:"_embedded"`
	}{}
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		log.Errorf("could not unmarshal deployments response: %#v", body)
		return nil, err
	}
	return response.Embedded.Deployments, nil
}

func (c *Client) GetDeployment(deploymentID string) (*Deployment, error) {
	body, err := c.Get(fmt.Sprintf("deployments/%s", deploymentID))
	if err != nil {
		log.Errorf("could not find Compose.io deployment %s: %s", deploymentID, err)
		return nil, err
	}

	deployment := &Deployment{}
	if err := json.Unmarshal([]byte(body), deployment); err != nil {
		log.Errorf("could not unmarshal deployment response: %#v", body)
		return nil, err
	}
	return deployment, nil
}

func (c *Client) GetDeploymentByName(name string) (*Deployment, error) {
	deployments, err := c.GetDeployments()
	if err != nil {
		return nil, err
	}

	for _, deployment := range deployments {
		if deployment.Name == name {
			return c.GetDeployment(deployment.ID)
		}
	}
	return nil, fmt.Errorf("could not find Compose.io deployment %s", name)
}

func (c *Client) DeleteDeployment(deploymentID string) (*Recipe, error) {
	body, err := c.Delete(fmt.Sprintf("deployments/%s", deploymentID))
	if err != nil {
		log.Errorf("could not delete Compose.io deployment %s: %s", deploymentID, err)
		return nil, err
	}

	recipe := &Recipe{}
	if err := json.Unmarshal([]byte(body), recipe); err != nil {
		log.Errorf("could not unmarshal recipe response: %#v", body)
		return nil, err
	}
	return recipe, nil
}
