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

func (c *Client) GetScalingByDeploymentID(id string) (*Scaling, error) {
	body, err := c.Get(fmt.Sprintf("deployments/%s/scalings", id))
	if err != nil {
		log.Errorf("could not get Compose.io scaling for deployment %s: %s", id, err)
		return nil, err
	}

	scaling := &Scaling{}
	if err := json.Unmarshal([]byte(body), scaling); err != nil {
		log.Errorf("could not unmarshal scaling response: %#v", body)
		return nil, err
	}
	return scaling, nil
}
