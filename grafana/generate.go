package grafana

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gracew/widget/graph/model"
	"github.com/pkg/errors"
)

const (
	CREATE = "create"
	READ   = "read"
	LIST   = "list"
	DELETE = "delete"
)

var (
	total    = "Total"
	method   = "{{method}}"
	quantile = "{{quantile}}"
)

func ImportDashboard(apiName string, deploy model.Deploy, customLogic *model.AllCustomLogic) error {
	req := createDashboardRequest{
		Dashboard: generateDashboard(apiName, deploy, customLogic),
		FolderID:  0,
		Overwrite: true,
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return errors.Wrap(err, "could not marshal grafana dashboard definition")
	}
	res, err := http.Post("http://admin:admin@localhost:3001/api/dashboards/db", "application/json", bytes.NewReader(reqBytes))
	if err != nil {
		return errors.Wrap(err, "could not create grafana dashboard")
	}
	resBytes, _ := ioutil.ReadAll(res.Body)
	log.Printf("grafana dashboard creation response: %s", string(resBytes))
	return nil
}

type createDashboardRequest struct {
	Dashboard Dashboard `json:"dashboard"`
	FolderID  int       `json:"folderId"`
	Overwrite bool      `json:"overwrite"`
}

type createDashboardResponse struct {
	UID    string `json:"uid"`
	URL    string `json:"url"`
	Status string `json:"status"`
}

type panelInput struct {
	title  string
	expr   string
	legend *string
}

func generateDashboard(apiName string, deploy model.Deploy, customLogic *model.AllCustomLogic) Dashboard {
	inputs := []panelInput{
		// overall
		panelInput{title: "Total Requests/sec (5 min avg)", expr: fmt.Sprintf("sum(rate(%s_http_requests_total[5m]))", apiName), legend: &total},
		panelInput{title: "Requests/sec by method (5 min avg)", expr: fmt.Sprintf("rate(%s_http_requests_total[5m])", apiName), legend: &method},
		// create
		panelInput{title: "Request Latency: Create", expr: fmt.Sprintf("%s_http_request_duration_seconds{method=\"%s\"}", apiName, CREATE), legend: &quantile},
	}
	if customLogic != nil {
		inputs = append(inputs, customLogicPanels(apiName, CREATE, customLogic.Create)...)
	}
	inputs = append(inputs,
		panelInput{title: "Database Latency: Create", expr: fmt.Sprintf("%s_database_access_duration_seconds{method=\"%s\"}", apiName, CREATE), legend: &quantile},
		// read
		panelInput{title: "Request Latency: Read", expr: fmt.Sprintf("%s_http_request_duration_seconds{method=\"%s\"}", apiName, READ), legend: &quantile},
		panelInput{title: "Database Latency: Read", expr: fmt.Sprintf("%s_database_access_duration_seconds{method=\"%s\"}", apiName, READ), legend: &quantile},
		// list
		panelInput{title: "Request Latency: List", expr: fmt.Sprintf("%s_http_request_duration_seconds{method=\"%s\"}", apiName, LIST), legend: &quantile},
		panelInput{title: "Database Latency: List", expr: fmt.Sprintf("%s_database_access_duration_seconds{method=\"%s\"}", apiName, LIST), legend: &quantile},
	)

	// update
	if customLogic != nil {
		for actionName, actionCustomLogic := range customLogic.Update {
			inputs = append(inputs, panelInput{title: "Request Latency: " + actionName, expr: fmt.Sprintf("%s_http_request_duration_seconds{method=\"%s\"}", apiName, actionName), legend: &quantile})
			inputs = append(inputs, customLogicPanels(apiName, actionName, actionCustomLogic)...)
			inputs = append(inputs, panelInput{title: "Database Latency: " + actionName, expr: fmt.Sprintf("%s_database_access_duration_seconds{method=\"%s\"}", apiName, actionName), legend: &quantile})
		}
	}

	// delete
	inputs = append(inputs, panelInput{title: "Request Latency: Delete", expr: fmt.Sprintf("%s_http_request_duration_seconds{method=\"%s\"}", apiName, DELETE), legend: &quantile})
	if customLogic != nil {
		inputs = append(inputs, customLogicPanels(apiName, DELETE, customLogic.Delete)...)
	}
	inputs = append(inputs, panelInput{title: "Database Latency: Delete", expr: fmt.Sprintf("%s_database_access_duration_seconds{method=\"%s\"}", apiName, DELETE), legend: &quantile})
	var panels []Panel
	for i, input := range inputs {
		panels = append(panels, generatePanel(apiName, i, input))
	}

	return Dashboard{
		Editable:      false,
		GraphTooltip:  2,
		Panels:        panels,
		Title:         apiName,
		SchemaVersion: 22,
		UID:           deploy.ID,
	}
}

func customLogicPanels(apiName string, method string, customLogic *model.CustomLogic) []panelInput {
	panels := []panelInput{}
	if customLogic == nil {
		return panels
	}
	if customLogic.Before != nil {
		panels = append(panels, panelInput{
			title:  fmt.Sprintf("Custom Logic Latency: before%s", strings.Title(strings.ToLower(method))),
			expr:   fmt.Sprintf("%s_custom_logic_duration_seconds{method=\"%s\", when=\"before\"}", apiName, method),
			legend: &quantile,
		})
	}
	if customLogic.After != nil {
		panels = append(panels, panelInput{
			title:  fmt.Sprintf("Custom Logic Latency: after%s", strings.Title(strings.ToLower(method))),
			expr:   fmt.Sprintf("%s_custom_logic_duration_seconds{method=\"%s\", when=\"after\"}", apiName, method),
			legend: &quantile,
		})
	}
	return panels
}

func generatePanel(apiName string, i int, input panelInput) Panel {
	var x int
	if i%2 == 0 {
		x = 0
	} else {
		x = 12
	}
	var legend string
	if input.legend != nil {
		legend = *input.legend
	}
	return Panel{
		ID:         i,
		Datasource: "Prometheus",
		GridPos:    GridPos{H: 9, W: 12, X: x, Y: i / 2 * 9},
		Targets: []Target{Target{
			Expr:         input.expr,
			LegendFormat: legend,
			RefID:        "A",
		}},
		Title: input.title,
		Type:  "graph",
	}
}
