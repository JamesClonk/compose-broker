package api

import (
	"encoding/json"

	"github.com/JamesClonk/compose-broker/log"
)

type Databases []Database
type Database struct {
	DatabaseType string `json:"type"`
	Status       string `json:"status"`
	Embedded     struct {
		Versions []Version `json:"versions"`
	} `json:"_embedded"`
}
type Version struct {
	Application string `json:"application"`
	Status      string `json:"status"`
	Preferred   bool   `json:"preferred"`
	Version     string `json:"version"`
}

func (c *Client) GetDatabases() (Databases, error) {
	body, err := c.GetJSON("databases")
	if err != nil {
		log.Errorf("could not get Compose.io databases: %s", err)
		return nil, err
	}

	type embeddedResponse struct {
		Embedded struct {
			Databases Databases `json:"applications"`
		} `json:"_embedded"`
	}
	response := embeddedResponse{}
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		log.Errorf("could not unmarshal databases response: %#v", body)
		return nil, err
	}
	return response.Embedded.Databases, nil
}
