// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
package graph

import (
	"context"

	"github.com/go-pg/pg"
	"github.com/gracew/widget/graph/model"
)

func (r *mutationResolver) SaveCustomLogic(ctx context.Context, input model.SaveCustomLogicInput) (bool, error) {
	// TODO(gracew): postgres probably isn't the best place for this
	db := pg.Connect(&pg.Options{User: "postgres"})
	defer db.Close()

	customLogic := &model.CustomLogic{
		APIID:         input.APIID,
		OperationType: input.OperationType,
		Language:      input.Language,
		Before:        input.Before,
		After:         input.After,
	}

	_, err := db.Model(customLogic).OnConflict("(apiid, operation_type) DO UPDATE").Insert()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *queryResolver) CustomLogic(ctx context.Context, apiID string) ([]*model.CustomLogic, error) {
	return r.Store.CustomLogic(apiID)
}
