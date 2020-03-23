// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
package graph

import (
	"context"

	"github.com/gracew/widget/graph/model"
)

func (r *mutationResolver) AddTestToken(ctx context.Context, input model.TestTokenInput) (*model.TestToken, error) {
	return r.Store.NewTestToken(input), nil
}

func (r *queryResolver) TestTokens(ctx context.Context) (*model.TestTokenResponse, error) {
	return &model.TestTokenResponse{TestTokens: r.Store.TestTokens()}, nil
}
