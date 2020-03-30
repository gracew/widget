package grafana

type Dashboard struct {
	Editable      bool    `json:"editable"`
	GraphTooltip  int     `json:"graphTooltip"`
	Panels        []Panel `json:"panels"`
	Title         string  `json:"title"`
	SchemaVersion int     `json:"schemaVersion"`
	UID           string  `json:"uid"`
}

type Panel struct {
	Datasource string   `json:"datasource"`
	GridPos    GridPos  `json:"gridPos"`
	ID         int      `json:"id"`
	Targets    []Target `json:"targets"`
	Title      string   `json:"title"`
	Type       string   `json:"type"`
}

type GridPos struct {
	H int `json:"h"`
	W int `json:"w"`
	X int `json:"x"`
	Y int `json:"y"`
}

type Target struct {
	Expr         string `json:"expr"`
	LegendFormat string `json:"legendFormat"`
	RefID        string `json:"refId"`
}
