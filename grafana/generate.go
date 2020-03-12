package grafana

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gracew/widget/graph/model"
	"github.com/pkg/errors"
)

func ImportDashboard(apiName string, deploy model.Deploy) error {
	req := createDashboardRequest{
		Dashboard: generateDashboard(apiName, deploy),
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

func generateDashboard(apiName string, deploy model.Deploy) Dashboard {
	panels := []Panel{
		// create
		generatePanel(apiName, 0, "Request Latency: Create", "http_request_duration_seconds{method=\"CREATE\"}"),
		generatePanel(apiName, 1, "Custom Logic Latency: beforeCreate", "custom_logic_duration_seconds{method=\"CREATE\", when=\"before\"}"),
		generatePanel(apiName, 2, "Custom Logic Latency: afterCreate", "custom_logic_duration_seconds{method=\"CREATE\", when=\"after\"}"),
		generatePanel(apiName, 3, "Database Latency: Create", "database_access_duration_seconds{method=\"CREATE\"}"),
		// read
		generatePanel(apiName, 4, "Request Latency: Read", "http_request_duration_seconds{method=\"READ\"}"),
		generatePanel(apiName, 5, "Database Latency: Read", "database_access_duration_seconds{method=\"READ\"}"),
		// list
		generatePanel(apiName, 6, "Request Latency: List", "http_request_duration_seconds{method=\"LIST\"}"),
		generatePanel(apiName, 7, "Database Latency: List", "database_access_duration_seconds{method=\"LIST\"}"),
	}

	return Dashboard{
		Editable: false,
		GraphTooltip: 1,
		Panels: panels,
		Title: apiName,
		SchemaVersion: 22,
		UID: deploy.ID,
	}
}

func generatePanel(apiName string, i int, title string, expr string) Panel {
	var x int
	if i % 2 == 0 {
		x = 0
	} else {
		x = 12
	}
	return Panel{
		ID: i,
		Datasource: "Prometheus",
		GridPos: GridPos{H: 9, W: 12, X: x, Y: i / 2 * 9},
		Targets: []Target{Target{
			Expr: apiName + "_" + expr,
			LegendFormat: "{{quantile}}",
			RefID: "A",
		}},
		Title: title,
		Type: "graph",
	}
}
