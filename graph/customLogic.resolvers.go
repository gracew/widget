// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
package graph

import (
	"context"
	"errors"

	"github.com/gracew/widget/graph/model"
)

func (r *actionDefinitionResolver) CustomLogic(ctx context.Context, obj *model.ActionDefinition) (*model.CustomLogic, error) {
	apiID := apiID(ctx)
	if apiID == nil {
		return nil, errors.New("expected API ID to be set in context")
	}
	customLogic, err := r.Store.CustomLogic(*apiID)
	if err != nil || customLogic == nil {
		return nil, err
	}
	return customLogic.Update[obj.Name], nil
}

func (r *createDefinitionResolver) CustomLogic(ctx context.Context, obj *model.CreateDefinition) (*model.CustomLogic, error) {
	apiID := apiID(ctx)
	if apiID == nil {
		return nil, errors.New("expected API ID to be set in context")
	}
	customLogic, err := r.Store.CustomLogic(*apiID)
	if err != nil || customLogic == nil {
		return nil, err
	}
	return customLogic.Create, nil
}

func (r *deleteDefinitionResolver) CustomLogic(ctx context.Context, obj *model.DeleteDefinition) (*model.CustomLogic, error) {
	apiID := apiID(ctx)
	if apiID == nil {
		return nil, errors.New("expected API ID to be set in context")
	}
	customLogic, err := r.Store.CustomLogic(*apiID)
	if err != nil || customLogic == nil {
		return nil, err
	}
	return customLogic.Delete, nil
}

func (r *mutationResolver) SaveCustomLogic(ctx context.Context, input model.SaveCustomLogicInput) (bool, error) {
	// TODO(gracew): postgres probably isn't the best place for this
	err := r.Store.SaveCustomLogic(input)
	if err != nil {
		return false, err
	}
	return true, nil
}
