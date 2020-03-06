package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gracew/widget/graph/model"
	"github.com/gracew/widget/parse"
	"github.com/gracew/widget/store"
)

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == http.MethodOptions {
		return
	}

	// get the userId
	parseToken := r.Header["X-Parse-Session-Token"][0]
	userID, err := parse.GetUserId(parseToken)
	if err != nil {
		panic(err)
	}

	// execute beforeSave logic
	vars := mux.Vars(r)
	apiID := vars["apiID"]

	customLogic, err := store.CustomLogicByOperation(apiID, model.OperationTypeCreate)
	if err != nil {
		panic(err)
	}

	originalReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	modifiedReqBytes := originalReq;
	if (customLogic != nil && customLogic.BeforeSave != nil) {
		modifiedReqBytes, err = Execute(originalReq, *customLogic.BeforeSave)
		if err != nil {
			panic(err)
		}
	}

	// add createdBy to the original req
	var modifiedReq map[string]interface{}
	err = json.Unmarshal(modifiedReqBytes, &modifiedReq)
	modifiedReq["createdBy"] = userID

	// delegate to parse
	res, err := parse.CreateObject(apiID, vars["env"], modifiedReq)
	if err != nil {
		panic(err)
	}

	// execute afterSave logic

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		panic(err)
	}
}

func ReadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == http.MethodOptions {
		return
	}

	// get the userId
	parseToken := r.Header["X-Parse-Session-Token"][0]
	userID, err := parse.GetUserId(parseToken)
	if err != nil {
		panic(err)
	}

	// delegate to parse
	vars := mux.Vars(r)
	res, err := parse.GetObject(vars["apiID"], vars["env"], vars["id"])
	if err != nil {
		panic(err)
	}

	// fetch the authorization policy
	// TODO(gracew): parallelize some of these requests
	auth, err := store.Auth(vars["apiID"])
	if err != nil {
		panic(err)
	}

	if auth.ReadPolicy.Type == model.AuthPolicyTypeCreatedBy {
		if userID != (*res)["createdBy"] {
			json.NewEncoder(w).Encode(&unauthorized{Message: "unauthorized"})
			return
		}
	}
	// TODO(gracew): support other authz policies

	json.NewEncoder(w).Encode(&res)
}

func ListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == http.MethodOptions {
		return
	}

	// get the userId
	parseToken := r.Header["X-Parse-Session-Token"][0]
	userID, err := parse.GetUserId(parseToken)
	if err != nil {
		panic(err)
	}

	// delegate to parse
	pageSizes, ok := r.URL.Query()["pageSize"]
	pageSize := "100"
	if ok && len(pageSizes[0]) >= 1 {
		pageSize = pageSizes[0]
	}
	vars := mux.Vars(r)
	res, err := parse.ListObjects(vars["apiID"], vars["env"], pageSize)
	if err != nil {
		panic(err)
	}

	// fetch the authorization policy
	// TODO(gracew): parallelize some of these requests
	auth, err := store.Auth(vars["apiID"])
	if err != nil {
		panic(err)
	}

	var filtered []parse.ObjectRes
	if auth.ReadPolicy.Type == model.AuthPolicyTypeCreatedBy {
		for i := 0; i < len(res.Results); i++ {
			if userID == res.Results[i]["createdBy"] {
				filtered = append(filtered, res.Results[i])
			}
		}
	}
	// TODO(gracew): support other authz policies

	json.NewEncoder(w).Encode(&parse.ListRes{Results: filtered})
}

type createdBy struct {
	CreatedBy string `json:"createdBy"`
}

type unauthorized struct {
	Message string `json:"message"`
}