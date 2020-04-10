package model

type API struct {
	ID         string               `json:"id"`
	Name       string               `json:"name"`
	Fields     []*FieldDefinition   `json:"fields"`
	Operations *OperationDefinition `json:"operations"`

	// nested resolver
	DeployIds []string `json:"deploys"`
}

type CreateDefinition struct {
	Enabled bool `json:"enabled"`

	CustomLogic string `json:"customLogic"`
}

type ReadDefinition struct {
	Enabled bool `json:"enabled"`

	Auth string `json:"auth"`
}

type ListDefinition struct {
	Enabled bool             `json:"enabled"`
	Sort    []SortDefinition `json:"sort"`
	Filter  []string         `json:"filter"`

	Auth string `json:"auth"`
}

type ActionDefinition struct {
	Name   string   `json:"name"`
	Fields []string `json:"fields"`

	Auth        string `json:"auth"`
	CustomLogic string `json:"customLogic"`
}

type DeleteDefinition struct {
	Enabled bool `json:"enabled"`

	Auth        string `json:"auth"`
	CustomLogic string `json:"customLogic"`
}

type Auth struct {
	APIID  string                 `json:"apiID" sql:",pk"`
	Read   *AuthPolicy            `json:"read"`
	Update map[string]*AuthPolicy `json:"update"`
	Delete *AuthPolicy            `json:"delete"`
}

type AllCustomLogic struct {
	APIID    string                  `json:"apiID" sql:",pk"`
	ImageURL string                  `json:"imageURL"`
	Create   *CustomLogic            `json:"create"`
	Update   map[string]*CustomLogic `json:"update"`
	Delete   *CustomLogic            `json:"delete"`
}

type DeployStepStatus struct {
	DeployID string       `json:"deployID" sql:",pk"`
	Step     DeployStep   `json:"step" sql:",pk"`
	Status   DeployStatus `json:"status"`
}
