// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
package graph

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/go-pg/pg"
	"github.com/google/uuid"
	"github.com/gracew/widget/graph/generated"
	"github.com/gracew/widget/graph/model"
	"github.com/gracew/widget/launch"
	"github.com/gracew/widget/store"
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

	currAPI := &model.API{ID: input.APIID}
	err := db.Select(currAPI)
	if err != nil {
		return nil, err
	}

	if definition.Name != currAPI.Name {
		return nil, errors.New("cannot change API name")
	}
	updatedAPI := &model.API{
		ID:         input.APIID,
		Definition: &definition,
	}
	err = db.Update(updatedAPI)
	if err != nil {
		return nil, err
	}
	return updatedAPI, nil
}

func (r *mutationResolver) AuthAPI(ctx context.Context, input model.AuthAPIInput) (bool, error) {
	db := pg.Connect(&pg.Options{User: "postgres"})
	defer db.Close()

	err := db.Insert(&model.Auth{
		ID:                 uuid.New().String(),
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

func (r *mutationResolver) DeployAPI(ctx context.Context, input model.DeployAPIInput) (*model.Deploy, error) {
	db := pg.Connect(&pg.Options{User: "postgres"})
	defer db.Close()

	auth, err := store.Auth(input.APIID)
	if err != nil {
		return nil, errors.Wrapf(err, "could not fetch auth for api %s", input.APIID)
	}

	deploy := &model.Deploy{
		ID:    uuid.New().String(),
		APIID: input.APIID,
		Env:   input.Env,
	}

	err = launch.DeployAPI(deploy.ID, *auth)
	if err != nil {
		return nil, errors.Wrapf(err, "could not launch container for api %s", input.APIID)
	}

	err = db.Insert(deploy)
	if err != nil {
		return nil, errors.Wrapf(err, "could not save deploy metadata for api %s", input.APIID)
	}

	return deploy, nil
}

func (r *mutationResolver) SaveCustomLogic(ctx context.Context, input model.SaveCustomLogicInput) (bool, error) {
	// TODO(gracew): postgres probably isn't the best place for this
	db := pg.Connect(&pg.Options{User: "postgres"})
	defer db.Close()

	customLogic := &model.CustomLogic{
		APIID:         input.APIID,
		OperationType: input.OperationType,
		BeforeSave:          input.BeforeSave,
		AfterSave:   input.AfterSave,
	}
	err := db.Insert(customLogic)
	if err != nil {
		return false, err
	}
	return true, nil
}

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

func (r *queryResolver) API(ctx context.Context, id string) (*model.API, error) {
	return store.API(id)
}

func (r *queryResolver) Apis(ctx context.Context) ([]*model.API, error) {
	return store.Apis()
}

func (r *queryResolver) Auth(ctx context.Context, apiID string) (*model.Auth, error) {
	return store.Auth(apiID)
}

func (r *queryResolver) CustomLogic(ctx context.Context, apiID string) ([]*model.CustomLogic, error) {
	return store.CustomLogic(apiID)
}

func (r *queryResolver) TestTokens(ctx context.Context) (*model.TestTokenResponse, error) {
	db := pg.Connect(&pg.Options{User: "postgres"})
	defer db.Close()

	var tokens []model.TestToken
	db.Model(&tokens).Select()

	var res []*model.TestToken
	for i := 0; i < len(tokens); i++ {
		res = append(res, &tokens[i])
	}
	return &model.TestTokenResponse{TestTokens: res}, nil
}

func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }
func (r *Resolver) Query() generated.QueryResolver       { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
