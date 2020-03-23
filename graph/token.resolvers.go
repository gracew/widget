// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
package graph

import (
	"context"

	"github.com/go-pg/pg"
	"github.com/gracew/widget/graph/model"
)

func (r *mutationResolver) AddTestToken(ctx context.Context, input model.TestTokenInput) (*model.TestToken, error) {
	db := pg.Connect(&pg.Options{User: "postgres"})
	defer db.Close()

	token := &model.TestToken{
		// TODO(gracew): enforce label uniqueness?
		Label: input.Label,
		Token: input.Token,
	}
	err := db.Insert(token)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (r *queryResolver) TestTokens(ctx context.Context) (*model.TestTokenResponse, error) {
	db := pg.Connect(&pg.Options{User: "postgres"})
	defer db.Close()

	var tokens []model.TestToken
	db.Model(&tokens).Select()

	return &model.TestTokenResponse{TestTokens: tokens}, nil
}
