package store

import (
	"github.com/go-pg/pg"
	"github.com/gracew/widget/graph/model"
)

// API fetches an API by ID.
func API(id string) (*model.API, error) {
	db := pg.Connect(&pg.Options{User: "postgres"})
	defer db.Close()

	api := &model.API{ID: id}
	err := db.Select(api)
	if err != nil {
		return nil, err
	}

	return api, nil
}

// Apis fetches all APIs.
func Apis() ([]*model.API, error) {
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

// Auth fetches auth for the specified API.
func Auth(apiID string) (*model.Auth, error) {
	db := pg.Connect(&pg.Options{User: "postgres"})
	defer db.Close()

	var auths []model.Auth
	err := db.Model(&auths).WhereIn("apiid IN (?)", []string{apiID}).Select()
	if err != nil {
		return nil, err
	}

	if len(auths) == 0 {
		return nil, nil
	}
	return &auths[0], nil
}

func CustomLogic(apiID string) ([]*model.CustomLogic, error) {
	db := pg.Connect(&pg.Options{User: "postgres"})
	defer db.Close()

	var models []model.CustomLogic
	err := db.Model(&models).WhereIn("apiid IN (?)", []string{apiID}).Select()
	if err != nil {
		return nil, err
	}

	var res []*model.CustomLogic
	for i := 0; i < len(models); i++ {
		res = append(res, &models[i])
	}

	return res, nil
}

func Deploys(apiID string) ([]*model.Deploy, error) {
	db := pg.Connect(&pg.Options{User: "postgres"})
	defer db.Close()

	var models []model.Deploy
	err := db.Model(&models).WhereIn("apiid IN (?)", []string{apiID}).Select()
	if err != nil {
		return nil, err
	}

	var res []*model.Deploy
	for i := 0; i < len(models); i++ {
		res = append(res, &models[i])
	}

	return res, nil
}
