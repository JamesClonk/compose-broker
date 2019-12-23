package api

import (
	"encoding/json"

	"github.com/JamesClonk/compose-broker/log"
)

type Accounts []Account
type Account struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func (c *Client) GetAccounts() (Accounts, error) {
	body, err := c.Get("accounts")
	if err != nil {
		log.Errorf("could not get Compose.io accounts: %s", err)
		return nil, err
	}

	response := struct {
		Embedded struct {
			Accounts Accounts `json:"accounts"`
		} `json:"_embedded"`
	}{}
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		log.Errorf("could not unmarshal accounts response: %#v", body)
		return nil, err
	}
	return response.Embedded.Accounts, nil
}
