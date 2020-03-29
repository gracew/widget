// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
package graph

import (
	"context"
	"errors"

	"github.com/gracew/widget/graph/model"
)

func (r *actionDefinitionResolver) Auth(ctx context.Context, obj *model.ActionDefinition) (*model.AuthPolicy, error) {
	apiID := apiID(ctx)
	if apiID == "" {
		return nil, errors.New("expected API ID to be set in context")
	}
	auth, err := r.Store.Auth(apiID)
	if err != nil || auth == nil {
		return nil, err
	}
	return auth.Update[obj.Name], nil
}

func (r *deleteDefinitionResolver) Auth(ctx context.Context, obj *model.DeleteDefinition) (*model.AuthPolicy, error) {
	apiID := apiID(ctx)
	if apiID == "" {
		return nil, errors.New("expected API ID to be set in context")
	}
	auth, err := r.Store.Auth(apiID)
	if err != nil || auth == nil {
		return nil, err
	}
	return auth.Delete, nil
}

func (r *listDefinitionResolver) Auth(ctx context.Context, obj *model.ListDefinition) (*model.AuthPolicy, error) {
	apiID := apiID(ctx)
	if apiID == "" {
		return nil, errors.New("expected API ID to be set in context")
	}
	auth, err := r.Store.Auth(apiID)
	if err != nil || auth == nil {
		return nil, err
	}
	return auth.Read, nil
}

func (r *mutationResolver) AuthAPI(ctx context.Context, input model.AuthAPIInput) (bool, error) {
	err := r.Store.SaveAuth(input)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *readDefinitionResolver) Auth(ctx context.Context, obj *model.ReadDefinition) (*model.AuthPolicy, error) {
	apiID := apiID(ctx)
	if apiID == "" {
		return nil, errors.New("expected API ID to be set in context")
	}
	auth, err := r.Store.Auth(apiID)
	if err != nil || auth == nil {
		return nil, err
	}
	return auth.Read, nil
}
