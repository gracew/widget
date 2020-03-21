package model

type API struct {
	ID         string               `json:"id"`
	Name       string               `json:"name"`
	Fields     []*FieldDefinition   `json:"fields"`
	DeployIds  []string      		`json:"deploys"`
	Operations *OperationDefinition `json:"operations"`
}

type OperationDefinition struct {
	APIID  string            `json:"apiID" sql:",pk"`
	Create *CreateDefinition `json:"create"`
	Read   *ReadDefinition   `json:"read"`
	List   *ListDefinition   `json:"list"`
}

type Auth struct {
	APIID              string             `json:"apiID" sql:",pk"`
	ReadPolicy         *AuthPolicy        `json:"readPolicy"`
	WritePolicy        *AuthPolicy        `json:"writePolicy"`
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
