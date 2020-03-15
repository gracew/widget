package generated

type Object struct {
	ID        string  `json:"id" sql:"type:uuid,default:gen_random_uuid()"`
	CreatedBy string  `json:"createdBy"`
	foo       bool    `json:"foo"`
	bar       float64 `json:"bar"`
	baz       int32   `json:"baz"`
	qux       string  `json:"qux"`
}
