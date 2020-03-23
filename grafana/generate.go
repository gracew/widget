package grafana

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gracew/widget/graph/model"
	"github.com/pkg/errors"
)

func ImportDashboard(apiName string, deploy model.Deploy, customLogic []model.CustomLogic) error {
	req := createDashboardRequest{
		Dashboard: generateDashboard(apiName, deploy, customLogic),
		FolderID: 0,
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
	FolderID int `json:"folderId"`
	Overwrite bool `json:"overwrite"`
}

type createDashboardResponse struct {
	UID string `json:"uid"`
	URL string `json:"url"`
	Status string `json:"status"`
}

type panelInput struct {
	title string
	expr string
	legend *string
}

func generateDashboard(apiName string, deploy model.Deploy, customLogic []model.CustomLogic) Dashboard {
	var createCustomLogic *model.CustomLogic
	for _, el := range customLogic {
		if el.OperationType == model.OperationTypeCreate {
			createCustomLogic = &el
		}
	}

	total := "Total"
	method := "{{method}}"
	quantile := "{{quantile}}"
	inputs := []panelInput{
		// overall
		panelInput{title: "Total Requests/sec (5 min avg)", expr: fmt.Sprintf("sum(rate(%s_http_requests_total[5m]))", apiName), legend: &total},
		panelInput{title: "Requests/sec by method (5 min avg)", expr: fmt.Sprintf("rate(%s_http_requests_total[5m])", apiName), legend: &method},
		// create
		panelInput{title: "Request Latency: Create", expr: apiName + "_http_request_duration_seconds{method=\"CREATE\"}", legend: &quantile},
	}
	if createCustomLogic != nil {
		if createCustomLogic.Before != nil {
			inputs = append(inputs, panelInput{
				title: "Custom Logic Latency: beforeCreate",
				expr: apiName + "_custom_logic_duration_seconds{method=\"CREATE\", when=\"before\"}",
				legend: &quantile,
			})
		}
		if createCustomLogic.After != nil {
			inputs = append(inputs, panelInput{
				title: "Custom Logic Latency: afterCreate",
				expr: apiName + "_custom_logic_duration_seconds{method=\"CREATE\", when=\"after\"}",
				legend: &quantile,
			})
		}
	}
	inputs = append(inputs,
		panelInput{title: "Database Latency: Create", expr: apiName + "_database_access_duration_seconds{method=\"CREATE\"}", legend: &quantile},
		// read
		panelInput{title: "Request Latency: Read", expr: apiName + "_http_request_duration_seconds{method=\"READ\"}", legend: &quantile},
		panelInput{title: "Database Latency: Read", expr: apiName + "_database_access_duration_seconds{method=\"READ\"}", legend: &quantile},
		// list
		panelInput{title: "Request Latency: List", expr: apiName + "_http_request_duration_seconds{method=\"LIST\"}", legend: &quantile},
		panelInput{title: "Database Latency: List", expr: apiName + "_database_access_duration_seconds{method=\"LIST\"}", legend: &quantile},
	)
	var panels []Panel
	for i, input := range inputs {
		panels = append(panels, generatePanel(apiName, i, input))
	}

	return Dashboard{
		Editable: false,
		GraphTooltip: 2,
		Panels: panels,
		Title: apiName,
		SchemaVersion: 22,
		UID: deploy.ID,
	}
}

func generatePanel(apiName string, i int, input panelInput) Panel {
	var x int
	if i % 2 == 0 {
		x = 0
	} else {
		x = 12
	}
	var legend string
	if input.legend != nil {
		legend = *input.legend
	}
	return Panel{
		ID: i,
		Datasource: "Prometheus",
		GridPos: GridPos{H: 9, W: 12, X: x, Y: i / 2 * 9},
		Targets: []Target{Target{
			Expr: input.expr,
			LegendFormat: legend,
			RefID: "A",
		}},
		Title: input.title,
		Type: "graph",
	}
}
