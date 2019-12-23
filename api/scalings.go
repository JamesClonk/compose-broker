package api

import (
	"encoding/json"
	"fmt"

	"github.com/JamesClonk/compose-broker/log"
)

type Scaling struct {
	AllocatedUnits     int    `json:"allocated_units"`
	UsedUnits          int    `json:"used_units"`
	StartingUnits      int    `json:"starting_units"`
	MinimumUnits       int    `json:"minimum_units"`
	MemoryPerUnitInMB  int    `json:"memory_per_unit_in_mb"`
	StoragePerUnitInMB int    `json:"storage_per_unit_in_mb"`
	UnitSizeInMB       int    `json:"unit_size_in_mb"`
	UnitType           string `json:"unit_type"`
}

func (c *Client) GetScaling(deploymentID string) (*Scaling, error) {
	body, err := c.Get(fmt.Sprintf("deployments/%s/scalings", deploymentID))
	if err != nil {
		log.Errorf("could not get Compose.io scaling for deployment %s: %s", deploymentID, err)
		return nil, err
	}

	scaling := &Scaling{}
	if err := json.Unmarshal([]byte(body), scaling); err != nil {
		log.Errorf("could not unmarshal scaling response: %#v", body)
		return nil, err
	}
	return scaling, nil
}

func (c *Client) UpdateScaling(deploymentID string, units int) (*Recipe, error) {
	body, err := c.Post(fmt.Sprintf("deployments/%s/scalings", deploymentID), fmt.Sprintf(`{"deployment":{"units":%d}}`, units))
	if err != nil {
		log.Errorf("could not update Compose.io scaling for deployment %s to %d units: %s", deploymentID, units, err)
		return nil, err
	}

	recipe := &Recipe{}
	if err := json.Unmarshal([]byte(body), recipe); err != nil {
		log.Errorf("could not unmarshal recipe response: %#v", body)
		return nil, err
	}
	return recipe, nil
}
