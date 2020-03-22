package store

import (
	"encoding/json"
	"log"

	"github.com/go-pg/pg"
	"github.com/google/uuid"
	"github.com/gracew/widget/graph/model"
	"github.com/pkg/errors"
)

type Store struct {
	DB *pg.DB
}

func (s Store) NewAPI(input model.DefineAPIInput) (*model.API, error) {
	// janky way of converting from DefineAPIInput -> API
	bytes, err := json.Marshal(input)
	if err != nil {
		return nil, errors.Wrapf(err, "could not marshal input to json");
	}

	var api model.API
	err = json.Unmarshal(bytes, &api)
	if err != nil {
		return nil, errors.Wrapf(err, "could not unmarshal bytes as json");
	}

	api.ID = uuid.New().String()
	err = s.DB.Insert(&api)
	if err != nil {
		return nil, err
	}
	return &api, nil
}

func (s Store) UpdateAPI(input model.UpdateAPIInput) (*model.API, error) {
	// janky way of converting from UpdateAPIInput -> API
	bytes, err := json.Marshal(input)
	if err != nil {
		return nil, errors.Wrapf(err, "could not marshal input to json");
	}

	var api model.API
	err = json.Unmarshal(bytes, &api)
	if err != nil {
		return nil, errors.Wrapf(err, "could not unmarshal bytes as json");
	}

	m := s.DB.Model(&api)
	if input.Fields != nil {
		m.Column("fields")
	}
	if input.Operations != nil {
		m.Column("operations")
	}
	// TODO(gracew): figure out better way to not clobber name
	_, err = m.WherePK().Update()
	if err != nil {
		return nil, err
	}
	return &api, nil
}

// API fetches an API by ID.
func (s Store) API(id string) (*model.API, error) {
	api := &model.API{ID: id}
	err := s.DB.Select(api)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch API %s", id)
	}

	return api, nil
}

// Apis fetches all APIs.
func (s Store) Apis() ([]*model.API, error) {
	var apis []model.API
	err := s.DB.Model(&apis).Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch APIs")
	}

	var res []*model.API
	for i := 0; i < len(apis); i++ {
		res = append(res, &apis[i])
	}
	return res, nil
}

func (s Store) DeleteApi(id string) error {
	api := &model.API{ID: id}
	err := s.DB.Delete(api)
	if err != nil {
		return errors.Wrapf(err, "failed to delete API %s", id)
	}
	return nil
}

func (s Store) SaveAuth(input model.AuthAPIInput) error {
	bytes, err := json.Marshal(input)
	if err != nil {
		return errors.Wrapf(err, "could not marshal input to json");
	}

	var auth model.Auth
	err = json.Unmarshal(bytes, &auth)
	if err != nil {
		return errors.Wrapf(err, "could not unmarshal bytes as json");
	}

	_, err = s.DB.Model(&auth).OnConflict("(apiid) DO UPDATE").Insert()
	if err != nil {
		return err
	}
	return nil
}

// Auth fetches auth for the specified API.
func (s Store) Auth(apiID string) (*model.Auth, error) {
	auth := &model.Auth{APIID: apiID}
	err := s.DB.Select(auth)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch auth for API %s", apiID)
	}

	return auth, nil
}

func (s Store) CustomLogic(apiID string) ([]*model.CustomLogic, error) {
	var models []model.CustomLogic
	err := s.DB.Model(&models).WhereIn("apiid IN (?)", []string{apiID}).Select()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch custom logic for API %s", apiID)
	}

	var res []*model.CustomLogic
	for i := 0; i < len(models); i++ {
		res = append(res, &models[i])
	}

	return res, nil
}

func (s Store) NewDeploy(deploy *model.Deploy) error {
	err := s.DB.Insert(deploy)
	if err != nil {
		return errors.Wrapf(err, "could not save deploy metadata for api %s", deploy.APIID)
	}
	return nil
}

func (s Store) DeleteDeploy(id string) error {
	err := s.DB.Delete(&model.Deploy{ID: id})
	if err != nil {
		return errors.Wrapf(err, "could not delete deploy %s", id)
	}
	return nil
}

func (s Store) Deploys(apiID string) ([]*model.Deploy, error) {
	var models []model.Deploy
	err := s.DB.Model(&models).WhereIn("apiid IN (?)", []string{apiID}).Select()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch deploys for API %s", apiID)
	}

	var res []*model.Deploy
	for i := 0; i < len(models); i++ {
		res = append(res, &models[i])
	}

	return res, nil
}

func (s Store) SaveDeployStepStatus(deployID string, step model.DeployStep, status model.DeployStatus) {
	stepStatus := &model.DeployStepStatus{
		DeployID: deployID,
		Step: step,
		Status: status,
	}
	_, err := s.DB.Model(stepStatus).OnConflict("(deploy_id, step) DO UPDATE").Insert()
	if err != nil {
		log.Printf("failed to record status for deploy %s, step %s, status %s: %v", deployID, step, status, err)
	}
}

func (s Store) DeployStatus(deployID string) ([]*model.DeployStepStatus, error) {
	var steps []model.DeployStepStatus
	err := s.DB.Model(&steps).WhereIn("deploy_id IN (?)", []string{deployID}).Select()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch step statuses for deploy %s", deployID)
	}

	var res []*model.DeployStepStatus
	for i := 0; i < len(steps); i++ {
		res = append(res, &steps[i])
	}

	return res, nil
}
