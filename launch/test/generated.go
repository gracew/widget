package generated

type Object struct {
	tableName struct{} `sql:"foo"`
	ID        string   `json:"id" sql:"type:uuid,default:gen_random_uuid()"`
	CreatedBy string   `json:"createdBy"`
	CreatedAt string   `json:"createdAt" sql:"default:now()"`
	Foo       bool     `json:"foo"`
	Bar       float64  `json:"bar"`
	Baz       int32    `json:"baz"`
	Qux       string   `json:"qux"`
	CamelCase string   `json:"camelCase"`
}
