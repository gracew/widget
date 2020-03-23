package store

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/gracew/widget/graph/model"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type Store struct {
	DB *gorm.DB
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
	s.DB.Create(&api)
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

	cols := []interface{}{}
	if input.Fields != nil {
		cols = append(cols, "fields")
	}
	if input.Operations != nil {
		cols = append(cols, "operations")
	}
	s.DB.Model(&api).Update(cols...)
	return &api, nil
}

// API fetches an API by ID.
func (s Store) API(id string) (*model.API, error) {
	api := &model.API{ID: id}
	s.DB.Select(api)
	return api, nil
}

// Apis fetches all APIs.
func (s Store) Apis() []*model.API {
	var apis []model.API
	s.DB.Find(&apis)

	var res []*model.API
	for i := 0; i < len(apis); i++ {
		res = append(res, &apis[i])
	}
	return res
}

func (s Store) DeleteApi(id string) {
	s.DB.Delete(&model.API{ID: id})
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

	s.DB.Save(&auth)
	return nil
}

// Auth fetches auth for the specified API.
func (s Store) Auth(apiID string) (*model.Auth, error) {
	auth := &model.Auth{APIID: apiID}
	s.DB.Select(auth)
	return auth, nil
}

func (s Store) SaveCustomLogic(input model.SaveCustomLogicInput) {
	s.DB.Save(&model.CustomLogic{
		APIID:         input.APIID,
		OperationType: input.OperationType,
		Language:      input.Language,
		Before:        input.Before,
		After:         input.After,
	})
}

func (s Store) CustomLogic(apiID string) []*model.CustomLogic {
	var models []model.CustomLogic
	s.DB.Where("apiid = ?", apiID).Find(&models)

	var res []*model.CustomLogic
	for i := 0; i < len(models); i++ {
		res = append(res, &models[i])
	}

	return res
}

func (s Store) NewDeploy(deploy *model.Deploy) {
	s.DB.Create(deploy)
}

func (s Store) DeleteDeploy(id string) {
	s.DB.Delete(&model.Deploy{ID: id})
}

func (s Store) Deploys(apiID string) []*model.Deploy {
	var models []model.Deploy
	s.DB.Where("apiid = ?", apiID).Find(&models)

	var res []*model.Deploy
	for i := 0; i < len(models); i++ {
		res = append(res, &models[i])
	}

	return res
}

func (s Store) SaveDeployStepStatus(deployID string, step model.DeployStep, status model.DeployStatus) {
	s.DB.Save(&model.DeployStepStatus{
		DeployID: deployID,
		Step: step,
		Status: status,
	})
}

func (s Store) DeployStatus(deployID string) []*model.DeployStepStatus {
	var steps []model.DeployStepStatus
	s.DB.Where("deploy_id = ?", deployID).Find(&steps)

	var res []*model.DeployStepStatus
	for i := 0; i < len(steps); i++ {
		res = append(res, &steps[i])
	}

	return res
}

func (s Store) NewTestToken(input model.TestTokenInput) *model.TestToken {
	token := &model.TestToken{
		// TODO(gracew): enforce label uniqueness?
		Label: input.Label,
		Token: input.Token,
	}
	s.DB.Create(token)
	return token
}

func (s Store) TestTokens() []*model.TestToken {
	var tokens []model.TestToken
	s.DB.Find(&tokens)

	var res []*model.TestToken
	for i := 0; i < len(tokens); i++ {
		res = append(res, &tokens[i])
	}
	return res
}
