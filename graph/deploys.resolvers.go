// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
package graph

import (
	"context"

	"github.com/go-pg/pg"
	"github.com/gracew/widget/grafana"
	"github.com/gracew/widget/graph/model"
	"github.com/gracew/widget/launch"
	"github.com/pkg/errors"
)

func (r *aPIResolver) Deploys(ctx context.Context, obj *model.API) ([]*model.Deploy, error) {
	return r.Store.Deploys(obj.ID)
}

func (r *mutationResolver) DeployAPI(ctx context.Context, input model.DeployAPIInput) (*model.Deploy, error) {
	db := pg.Connect(&pg.Options{User: "postgres"})
	defer db.Close()

	// TODO(gracew): parallelize these db calls lol
	api, err := r.Store.API(input.APIID)
	if err != nil {
		return nil, errors.Wrapf(err, "could not fetch api %s", input.APIID)
	}

	auth, err := r.Store.Auth(input.APIID)
	if err != nil {
		return nil, errors.Wrapf(err, "could not fetch auth for api %s", input.APIID)
	}

	customLogic, err := r.Store.CustomLogic(input.APIID)
	if err != nil {
		return nil, errors.Wrapf(err, "could not fetch custom logic for api %s", input.APIID)
	}

	deploy := &model.Deploy{
		ID:    input.DeployID,
		APIID: input.APIID,
		Env:   input.Env,
	}
	err = db.Insert(deploy)
	if err != nil {
		return nil, errors.Wrapf(err, "could not save deploy metadata for api %s", input.APIID)
	}

	launcher := launch.Launcher{
		Store:       r.Store,
		DeployID:    deploy.ID,
		API:         *api,
		Auth:        *auth,
		CustomLogic: customLogic,
	}

	err = launcher.DeployAPI()
	if err != nil {
		return nil, errors.Wrapf(err, "could not launch container for api %s", input.APIID)
	}

	if len(customLogic) > 0 {
		err = launcher.DeployCustomLogic()
		if err != nil {
			return nil, errors.Wrapf(err, "could not launch custom logic for api %s", input.APIID)
		}
	}

	err = grafana.ImportDashboard(api.Name, *deploy, customLogic)
	if err != nil {
		return nil, errors.Wrapf(err, "could not create grafana dashboard for api %s", input.APIID)
	}

	return deploy, nil
}

func (r *queryResolver) DeployStatus(ctx context.Context, deployID string) (*model.DeployStatusResponse, error) {
	steps, err := r.Store.DeployStatus(deployID)
	if err != nil {
		return nil, err
	}
	return &model.DeployStatusResponse{Steps: steps}, nil
}