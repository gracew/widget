package model

type API struct {
	ID         string         `json:"id"`
	Name       string         `json:"name"`
	DeployIds    []string      `json:"deploys"`
	Definition *APIDefinition `json:"definition"`
}