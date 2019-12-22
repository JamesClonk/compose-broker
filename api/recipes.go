package api

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/JamesClonk/compose-broker/log"
)

type Recipes []Recipe
type Recipe struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	Template           string    `json:"template"`
	Status             string    `json:"status"` // can be running, waiting, complete or failed
	StatusDetail       string    `json:"status_detail"`
	AccountID          string    `json:"account_id"`
	DeploymentID       string    `json:"deployment_id"`
	ParentID           string    `json:"parent_id"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	OperationsComplete int       `json:"operations_complete"`
	OperationsTotal    int       `json:"operations_total"`
	Embedded           struct {
		Recipes []Recipe `json:"recipes"`
	} `json:"_embedded"`
}

func (r Recipes) SortByCreatedAt() {
	sort.Slice(r, func(i, j int) bool {
		return r[j].CreatedAt.Before(r[i].CreatedAt)
	})
}

func (r Recipes) SortByUpdatedAt() {
	sort.Slice(r, func(i, j int) bool {
		return r[j].UpdatedAt.Before(r[i].UpdatedAt)
	})
}

func (c *Client) GetRecipe(id string) (*Recipe, error) {
	body, err := c.Get(fmt.Sprintf("recipes/%s", id))
	if err != nil {
		log.Errorf("could not get Compose.io recipe %s: %s", id, err)
		return nil, err
	}

	recipe := &Recipe{}
	if err := json.Unmarshal([]byte(body), recipe); err != nil {
		log.Errorf("could not unmarshal recipe response: %#v", body)
		return nil, err
	}
	return recipe, nil
}

func (c *Client) GetRecipesByDeploymentID(id string) (Recipes, error) {
	body, err := c.Get(fmt.Sprintf("deployments/%s/recipes", id))
	if err != nil {
		log.Errorf("could not get Compose.io recipes for deployment %s: %s", id, err)
		return nil, err
	}

	type embeddedResponse struct {
		Embedded struct {
			Recipes Recipes `json:"recipes"`
		} `json:"_embedded"`
	}
	response := embeddedResponse{}
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		log.Errorf("could not unmarshal recipes response: %#v", body)
		return nil, err
	}

	recipes := response.Embedded.Recipes
	recipes.SortByUpdatedAt()
	return recipes, nil
}
