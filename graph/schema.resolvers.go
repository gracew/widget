// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
package graph

import (
	"context"

	"github.com/go-pg/pg"
	"github.com/google/uuid"
	"github.com/gracew/widget/graph/generated"
	"github.com/gracew/widget/graph/model"
)

func (r *mutationResolver) CreateAPI(ctx context.Context, input model.NewAPI) (*model.API, error) {
	db := pg.Connect(&pg.Options{User: "postgres"})
	defer db.Close()

	api := &model.API{
		ID:   uuid.New().String(),
		Name: input.Name,
	}
	err := db.Insert(api)
	if err != nil {
		panic(err)
	}
	return api, nil
}

func (r *queryResolver) Apis(ctx context.Context) ([]*model.API, error) {
	api := model.API{ID: "1", Name: "comments"}
	return []*model.API{&api}, nil
}

func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }
func (r *Resolver) Query() generated.QueryResolver       { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
