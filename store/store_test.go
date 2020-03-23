package store

import (
	"testing"

	"github.com/google/uuid"
	"github.com/gracew/widget/graph/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type StoreTestSuite struct {
    suite.Suite
    store Store
}

func (suite *StoreTestSuite) SetupTest() {
	db, err := gorm.Open("postgres", "user=postgres sslmode=disable")
	if err != nil {
		panic(err)
	}
    suite.store = Store{DB: db}
}

func (suite *StoreTestSuite) TeardownTest() {
    suite.store.DB.Close()
}

func TestStoreTestSuite(t *testing.T) {
    suite.Run(t, new(StoreTestSuite))
}

func (s *StoreTestSuite) TestNewAPI() {
	name := uuid.New().String()
	input := model.DefineAPIInput{
		Name: name,
			Fields: []*model.FieldDefinitionInput{
			&model.FieldDefinitionInput{Name: "foo", Type: model.TypeBoolean,
		},
	}}
	createRes, err := s.store.NewAPI(input)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), name, createRes.Name)
	assert.Equal(s.T(), []*model.FieldDefinition{
		&model.FieldDefinition{Name: "foo", Type: model.TypeBoolean},
	}, createRes.Fields)

	getRes, err := s.store.API(createRes.ID)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), createRes, getRes)

	update := model.UpdateAPIInput{
		ID: createRes.ID,
		Fields: []*model.FieldDefinitionInput{
			&model.FieldDefinitionInput{Name: "foo", Type: model.TypeBoolean},
			&model.FieldDefinitionInput{Name: "bar", Type: model.TypeFloat},
		},
	}
	updateRes, err := s.store.UpdateAPI(update)
	assert.Nil(s.T(), err)
	// assert.Equal(s.T(), name, updateRes.Name)
	assert.Equal(s.T(), []*model.FieldDefinition{
		&model.FieldDefinition{Name: "foo", Type: model.TypeBoolean},
		&model.FieldDefinition{Name: "bar", Type: model.TypeFloat},
	}, updateRes.Fields)
}

func (s *StoreTestSuite) TestAuth() {
	apiID := uuid.New().String()
	input := model.AuthAPIInput{
		APIID: apiID,
		ReadPolicy: &model.AuthPolicyInput{Type: model.AuthPolicyTypeCreatedBy},
	}
	err := s.store.SaveAuth(input)
	assert.Nil(s.T(), err)

	auth, err := s.store.Auth(input.APIID)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), input.APIID, auth.APIID)
	assert.Equal(s.T(), &model.AuthPolicy{Type: model.AuthPolicyTypeCreatedBy}, auth.ReadPolicy)

	update := model.AuthAPIInput{
		APIID: apiID,
		ReadPolicy: &model.AuthPolicyInput{Type: model.AuthPolicyTypeAttributeMatch},
	}

	err = s.store.SaveAuth(update)
	assert.Nil(s.T(), err)

	auth2, err := s.store.Auth(input.APIID)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), &model.AuthPolicy{Type: model.AuthPolicyTypeAttributeMatch}, auth2.ReadPolicy)
}

func (s *StoreTestSuite) TestDeployStatus() {
	deployID := uuid.New().String()
	status := s.store.DeployStatus(deployID)
	assert.Nil(s.T(), status)

	s.store.SaveDeployStepStatus(deployID, model.DeployStepGenerateCode, model.DeployStatusComplete)
	s.store.SaveDeployStepStatus(deployID, model.DeployStepBuildImage, model.DeployStatusInProgress)
	status = s.store.DeployStatus(deployID)
	assert.Equal(s.T(), []*model.DeployStepStatus{
		&model.DeployStepStatus{DeployID: deployID, Step: model.DeployStepGenerateCode, Status: model.DeployStatusComplete},
		&model.DeployStepStatus{DeployID: deployID, Step: model.DeployStepBuildImage, Status: model.DeployStatusInProgress},
	}, status)

	s.store.SaveDeployStepStatus(deployID, model.DeployStepBuildImage, model.DeployStatusComplete)
	status = s.store.DeployStatus(deployID)
	assert.Equal(s.T(), []*model.DeployStepStatus{
		&model.DeployStepStatus{DeployID: deployID, Step: model.DeployStepGenerateCode, Status: model.DeployStatusComplete},
		&model.DeployStepStatus{DeployID: deployID, Step: model.DeployStepBuildImage, Status: model.DeployStatusComplete},
	}, status)
}

