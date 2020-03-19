// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
package graph

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/go-pg/pg"
	"github.com/google/uuid"
	"github.com/gracew/widget/graph/generated"
	"github.com/gracew/widget/graph/model"
)

func (r *mutationResolver) DefineAPI(ctx context.Context, input model.DefineAPIInput) (*model.API, error) {
	var definition model.APIDefinition
	json.Unmarshal([]byte(input.RawDefinition), &definition)

	db := pg.Connect(&pg.Options{User: "postgres"})
	defer db.Close()

	api := &model.API{
		ID:         uuid.New().String(),
		Name:       definition.Name,
		Definition: &definition,
	}
	err := db.Insert(api)
	if err != nil {
		return nil, err
	}
	return api, nil
}

func (r *mutationResolver) UpdateAPI(ctx context.Context, input model.UpdateAPIInput) (*model.API, error) {
	// validate input
	var definition model.APIDefinition
	json.Unmarshal([]byte(input.RawDefinition), &definition)

	db := pg.Connect(&pg.Options{User: "postgres"})
	defer db.Close()

	currAPI := &model.API{ID: input.ID}
	err := db.Select(currAPI)
	if err != nil {
		return nil, err
	}

	if definition.Name != currAPI.Name {
		return nil, errors.New("cannot change API name")
	}
	updatedAPI := &model.API{
		ID:         input.ID,
		Name:       currAPI.Name,
		Definition: &definition,
	}
	err = db.Update(updatedAPI)
	if err != nil {
		return nil, err
	}
	return updatedAPI, nil
}

func (r *mutationResolver) DeleteAPI(ctx context.Context, id string) (bool, error) {
	err := r.Store.DeleteApi(id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *queryResolver) API(ctx context.Context, id string) (*model.API, error) {
	return r.Store.API(id)
}

func (r *queryResolver) Apis(ctx context.Context) ([]*model.API, error) {
	return r.Store.Apis()
}

func (r *Resolver) API() generated.APIResolver           { return &aPIResolver{r} }
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }
func (r *Resolver) Query() generated.QueryResolver       { return &queryResolver{r} }

type aPIResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
