package model

type API struct {
	ID         string         `json:"id"`
	Name       string         `json:"name"`
	DeployIds    []string      `json:"deploys"`
	Definition *APIDefinition `json:"definition"`
}

type CustomLogic struct {
	APIID         string        `sql:",pk" json:"apiID"`
	OperationType OperationType `sql:",pk" json:"operationType"`
	Language      Language      `json:"language"`
	Before        *string       `json:"before"`
	After         *string       `json:"after"`
}