package pgmodel

import "github.com/gracew/widget/graph/model"

type CustomLogic struct {
	APIID         string        `sql:",pk" json:"apiID"`
	OperationType model.OperationType `sql:",pk" json:"operationType"`
	Language      model.Language      `json:"language"`
	Before        *string       `json:"before"`
	After         *string       `json:"after"`
}