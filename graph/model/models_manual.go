package model

type API struct {
	ID         string               `json:"id"`
	Name       string               `json:"name"`
	Fields     []*FieldDefinition   `json:"fields" gorm:"jsonb"`
	DeployIds  []string      		`json:"deploys"`
	Operations *OperationDefinition `json:"operations"`
}

type Auth struct {
	APIID              string             `json:"apiID" sql:",pk"`
	ReadPolicy         *AuthPolicy        `json:"readPolicy"`
	WritePolicy        *AuthPolicy        `json:"writePolicy"`
	DeletePolicy       *AuthPolicy        `json:"deletePolicy"`
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
