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
		return nil, err
	}
	return api, nil
}

func (r *mutationResolver) DeployAPI(ctx context.Context, input model.DeployAPI) (*model.Deploy, error) {
	db := pg.Connect(&pg.Options{User: "postgres"})
	defer db.Close()

	deploy := &model.Deploy{
		ID:    uuid.New().String(),
		APIID: input.APIID,
		Env:   input.Env,
	}
	// TODO(gracew): actually do some k8s/dns stuff lol
	err := db.Insert(deploy)
	if err != nil {
		return nil, err
	}
	return deploy, nil
}

func (r *queryResolver) Apis(ctx context.Context) ([]*model.API, error) {
	db := pg.Connect(&pg.Options{User: "postgres"})
	defer db.Close()

	var apis []model.API
	db.Model(&apis).Select()

	var res []*model.API
	for i := 0; i < len(apis); i++ {
		res = append(res, &apis[i])
	}
	return res, nil
}

func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }
func (r *Resolver) Query() generated.QueryResolver       { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
