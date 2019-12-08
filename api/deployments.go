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
	CreatedAt           time.Time `json:"created_at"`
	ConnectionStrings   struct {
		Direct []string `json:"direct"`
		CLI    []string `json:"cli"`
		Maps   []string `json:"maps"`
		SSH    []string `json:"ssh"`
		Health []string `json:"health"`
		Admin  []string `json:"admin"`
		Misc   []string `json:"misc"`
	} `json:"connection_strings"`
	Links struct {
		ComposeWebUI struct {
			HREF      string `json:"href"`
			Templated bool   `json:"templated"`
		} `json:"compose_web_ui"`
	} `json:"_links"`
}

func (c *Client) GetDeployments() (Deployments, error) {
	body, err := c.GetJSON("deployments")
	if err != nil {
		log.Errorf("could not get Compose.io deployments: %s", err)
		return nil, err
	}

	type embeddedResponse struct {
		Embedded struct {
			Deployments Deployments `json:"deployments"`
		} `json:"_embedded"`
	}
	response := embeddedResponse{}
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		log.Errorf("could not unmarshal deployments response: %#v", body)
		return nil, err
	}
	return response.Embedded.Deployments, nil
}

func (c *Client) GetDeploymentByID(id string) (*Deployment, error) {
	body, err := c.GetJSON(fmt.Sprintf("deployments/%s", id))
	if err != nil {
		log.Errorf("could not find Compose.io deployment %s: %s", id, err)
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
			return c.GetDeploymentByID(deployment.ID)
		}
	}
	return nil, fmt.Errorf("could not find Compose.io deployment %s", name)
}
