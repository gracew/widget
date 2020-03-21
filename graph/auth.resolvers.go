// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
package graph

import (
	"context"

	"github.com/gracew/widget/graph/model"
)

func (r *mutationResolver) AuthAPI(ctx context.Context, input model.AuthAPIInput) (bool, error) {
	err := r.Store.SaveAuth(input)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *queryResolver) Auth(ctx context.Context, apiID string) (*model.Auth, error) {
	return r.Store.Auth(apiID)
}
