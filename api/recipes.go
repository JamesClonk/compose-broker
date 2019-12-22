package api

import (
	"time"
)

type Recipes []Recipe
type Recipe struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	Template           string    `json:"template"`
	Status             string    `json:"status"`
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
