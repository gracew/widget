package parse

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type parseRes struct {
	CreatedAt string `json:"createdAt"`
	ObjectID  string `json:"objectId"`
}

type createdBy struct {
	CreatedBy string `json:"createdBy`
}

func GetUserId(parseToken string) (string, error) {
	req, err := http.NewRequest("GET", "http://localhost:1337/parse/users/me", nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("X-Parse-Application-Id", "appId")
	req.Header.Add("X-Parse-Session-Token", parseToken)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	var parseRes parseRes
	err = json.NewDecoder(res.Body).Decode(&parseRes)
	if err != nil {
		return "", err
	}
	return parseRes.ObjectID, nil
}

func CreateObject(apiID string, env string, req map[string]interface{}) ([]byte, error) {
	// TODO(gracew): don't hardcode this
	marshaled, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}
	parseReq, err := http.NewRequest("POST", "http://localhost:1337/parse/classes/"+parseClassName(apiID, env), bytes.NewReader(marshaled))
	if err != nil {
		panic(err)
	}
	parseReq.Header.Add("X-Parse-Application-Id", "appId")
	parseReq.Header.Add("Content-type", "application/json")
	client := &http.Client{}
	parseRes, err := client.Do(parseReq)
	if err != nil {
		panic(err)
	}

	return ioutil.ReadAll(parseRes.Body)
}

func GetObject(apiID string, env string, objectID string) ([]byte, error) {
	req, err := http.NewRequest("GET", "http://localhost:1337/parse/classes/"+parseClassName(apiID, env)+"/"+objectID, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("X-Parse-Application-Id", "appId")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	return ioutil.ReadAll(res.Body)
}

func ListObjects(apiID string, env string, pageSize string) ([]byte, error) {
	data := "limit=" + pageSize
	req, err := http.NewRequest("GET", "http://localhost:1337/parse/classes/"+parseClassName(apiID, env), strings.NewReader(data))
	if err != nil {
		panic(err)
	}
	req.Header.Add("X-Parse-Application-Id", "appId")
	req.Header.Add("Content-type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	return ioutil.ReadAll(res.Body)
}

func parseClassName(apiID string, env string) string {
	// parse class names cannot start with numbers or contain dashes
	return fmt.Sprintf("w%s_%s", strings.Replace(apiID, "-", "", -1), env)

}
