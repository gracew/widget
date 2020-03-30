package store

import (
	"testing"

	"github.com/go-pg/pg"
	"github.com/google/uuid"
	"github.com/gracew/widget/graph/model"
	"github.com/stretchr/testify/assert"
)

func TestNewAPI(t *testing.T) {
	db := pg.Connect(&pg.Options{User: "postgres"})
	defer db.Close()
	s := Store{DB: db}

	name := uuid.New().String()
	input := model.DefineAPIInput{
		Name: name,
		Fields: []model.FieldDefinitionInput{
			model.FieldDefinitionInput{Name: "foo", Type: model.TypeBoolean},
		}}
	createRes, err := s.NewAPI(input)
	assert.Nil(t, err)
	assert.Equal(t, name, createRes.Name)
	assert.Equal(t, []*model.FieldDefinition{
		&model.FieldDefinition{Name: "foo", Type: model.TypeBoolean},
	}, createRes.Fields)

	getRes, err := s.API(createRes.ID)
	assert.Nil(t, err)
	assert.Equal(t, createRes, getRes)

	update := model.UpdateAPIInput{
		ID: createRes.ID,
		Fields: []model.FieldDefinitionInput{
			model.FieldDefinitionInput{Name: "foo", Type: model.TypeBoolean},
			model.FieldDefinitionInput{Name: "bar", Type: model.TypeFloat},
		},
	}
	updateRes, err := s.UpdateAPI(update)
	assert.Nil(t, err)
	// assert.Equal(t, name, updateRes.Name)
	assert.Equal(t, []model.FieldDefinition{
		model.FieldDefinition{Name: "foo", Type: model.TypeBoolean},
		model.FieldDefinition{Name: "bar", Type: model.TypeFloat},
	}, updateRes.Fields)
}

func TestAuth(t *testing.T) {
	db := pg.Connect(&pg.Options{User: "postgres"})
	defer db.Close()
	s := Store{DB: db}

	apiID := uuid.New().String()
	input := model.AuthAPIInput{
		APIID:      apiID,
		ReadPolicy: &model.AuthPolicyInput{Type: model.AuthPolicyTypeCreatedBy},
	}
	err := s.SaveAuth(input)
	assert.Nil(t, err)

	auth, err := s.Auth(input.APIID)
	assert.Nil(t, err)
	assert.Equal(t, input.APIID, auth.APIID)
	assert.Equal(t, &model.AuthPolicy{Type: model.AuthPolicyTypeCreatedBy}, auth.ReadPolicy)

	update := model.AuthAPIInput{
		APIID:      apiID,
		ReadPolicy: &model.AuthPolicyInput{Type: model.AuthPolicyTypeAttributeMatch},
	}

	err = s.SaveAuth(update)
	assert.Nil(t, err)

	auth2, err := s.Auth(input.APIID)
	assert.Nil(t, err)
	assert.Equal(t, &model.AuthPolicy{Type: model.AuthPolicyTypeAttributeMatch}, auth2.ReadPolicy)
}

func TestDeployStatus(t *testing.T) {
	db := pg.Connect(&pg.Options{User: "postgres"})
	defer db.Close()
	s := Store{DB: db}

	deployID := uuid.New().String()
	status, err := s.DeployStatus(deployID)
	assert.NoError(t, err)
	assert.Nil(t, status)

	s.SaveDeployStepStatus(deployID, model.DeployStepGenerateCode, model.DeployStatusComplete)
	s.SaveDeployStepStatus(deployID, model.DeployStepBuildImage, model.DeployStatusInProgress)
	status, err = s.DeployStatus(deployID)
	assert.NoError(t, err)
	assert.Equal(t, []*model.DeployStepStatus{
		&model.DeployStepStatus{DeployID: deployID, Step: model.DeployStepGenerateCode, Status: model.DeployStatusComplete},
		&model.DeployStepStatus{DeployID: deployID, Step: model.DeployStepBuildImage, Status: model.DeployStatusInProgress},
	}, status)

	s.SaveDeployStepStatus(deployID, model.DeployStepBuildImage, model.DeployStatusComplete)
	status, err = s.DeployStatus(deployID)
	assert.NoError(t, err)
	assert.Equal(t, []*model.DeployStepStatus{
		&model.DeployStepStatus{DeployID: deployID, Step: model.DeployStepGenerateCode, Status: model.DeployStatusComplete},
		&model.DeployStepStatus{DeployID: deployID, Step: model.DeployStepBuildImage, Status: model.DeployStatusComplete},
	}, status)
}
