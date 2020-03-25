// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
package graph

import (
	"context"

	"github.com/gracew/widget/graph/generated"
	"github.com/gracew/widget/graph/model"
)

func (r *mutationResolver) DefineAPI(ctx context.Context, input model.DefineAPIInput) (*model.API, error) {
	return r.Store.NewAPI(input)
}

func (r *mutationResolver) UpdateAPI(ctx context.Context, input model.UpdateAPIInput) (*model.API, error) {
	return r.Store.UpdateAPI(input)
}

func (r *mutationResolver) DeleteAPI(ctx context.Context, id string) (bool, error) {
	err := r.Store.DeleteApi(id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *queryResolver) API(ctx context.Context, id string) (*model.API, error) {
	context.WithValue(ctx, apiIDCtxKey, id)
	return r.Store.API(id)
}

func (r *queryResolver) Apis(ctx context.Context) ([]model.API, error) {
	return r.Store.Apis()
}

func (r *Resolver) API() generated.APIResolver           { return &aPIResolver{r} }
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }
func (r *Resolver) Query() generated.QueryResolver       { return &queryResolver{r} }

type aPIResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
