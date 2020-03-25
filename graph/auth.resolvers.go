// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
package graph

import (
	"context"
	"fmt"

	"github.com/gracew/widget/graph/model"
)

func (r *actionDefinitionResolver) Auth(ctx context.Context, obj *model.ActionDefinition) (*model.AuthPolicy, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *deleteDefinitionResolver) Auth(ctx context.Context, obj *model.DeleteDefinition) (*model.AuthPolicy, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *listDefinitionResolver) Auth(ctx context.Context, obj *model.ListDefinition) (*model.AuthPolicy, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AuthAPI(ctx context.Context, input model.AuthAPIInput) (bool, error) {
	err := r.Store.SaveAuth(input)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *readDefinitionResolver) Auth(ctx context.Context, obj *model.ReadDefinition) (*model.AuthPolicy, error) {
	panic(fmt.Errorf("not implemented"))
}

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *queryResolver) Auth(ctx context.Context, apiID string) (*model.Auth, error) {
	return r.Store.Auth(apiID)
}
