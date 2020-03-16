package store

import (
	"testing"

	"github.com/go-pg/pg"
	"github.com/google/uuid"
	"github.com/gracew/widget/graph/model"
	"github.com/stretchr/testify/assert"
)

func TestDeployStatus(t *testing.T) {
	db := pg.Connect(&pg.Options{User: "postgres" })
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
