// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
package graph

import (
	"context"
	"fmt"

	"github.com/gracew/widget/graph/generated"
	"github.com/gracew/widget/graph/model"
)

func (r *mutationResolver) CreateAPI(ctx context.Context, input model.NewAPI) (*model.API, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Apis(ctx context.Context) ([]*model.API, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }
func (r *Resolver) Query() generated.QueryResolver       { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *mutationResolver) CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
	panic(fmt.Errorf("not implemented"))
}
func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
	u := model.User{ID: "1", Name: "Grace"}
	todo := model.Todo{ID: "1", Text: "Buy groceries", Done: true, User: &u}
	return []*model.Todo{&todo}, nil
}
