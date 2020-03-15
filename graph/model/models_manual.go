package model

type API struct {
	ID         string         `json:"id"`
	Name       string         `json:"name"`
	DeployIds    []string      `json:"deploys"`
	Definition *APIDefinition `json:"definition"`
}

type CustomLogic struct {
	APIID         string        `json:"apiID" sql:",pk"`
	OperationType OperationType `json:"operationType" sql:",pk"`
	Language      Language      `json:"language"`
	Before        *string       `json:"before"`
	After         *string       `json:"after"`
}

type DeployStepStatus struct {
	DeployID string       `json:"deployID" sql:",pk"`
	Step     DeployStep   `json:"step" sql:",pk"`
	Status   DeployStatus `json:"status"`
}
