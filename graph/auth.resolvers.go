// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
package graph

import (
	"context"

	"github.com/go-pg/pg"
	"github.com/gracew/widget/graph/model"
)

func (r *mutationResolver) AuthAPI(ctx context.Context, input model.AuthAPIInput) (bool, error) {
	db := pg.Connect(&pg.Options{User: "postgres"})
	defer db.Close()

	err := db.Insert(&model.Auth{
		APIID:              input.APIID,
		AuthenticationType: input.AuthenticationType,
		ReadPolicy: &model.AuthPolicy{
			Type:            input.ReadPolicy.Type,
			UserAttribute:   input.ReadPolicy.UserAttribute,
			ObjectAttribute: input.ReadPolicy.ObjectAttribute,
		},
		WritePolicy: &model.AuthPolicy{
			Type:            input.WritePolicy.Type,
			UserAttribute:   input.WritePolicy.UserAttribute,
			ObjectAttribute: input.WritePolicy.ObjectAttribute,
		},
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *queryResolver) Auth(ctx context.Context, apiID string) (*model.Auth, error) {
	return r.Store.Auth(apiID)
}
