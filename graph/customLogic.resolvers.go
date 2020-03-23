// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
package graph

import (
	"context"

	"github.com/gracew/widget/graph/model"
)

func (r *mutationResolver) SaveCustomLogic(ctx context.Context, input model.SaveCustomLogicInput) (bool, error) {
	// TODO(gracew): postgres probably isn't the best place for this
	r.Store.SaveCustomLogic(input)
	return true, nil
}

func (r *queryResolver) CustomLogic(ctx context.Context, apiID string) ([]*model.CustomLogic, error) {
	return r.Store.CustomLogic(apiID), nil
}
