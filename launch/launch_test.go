package launch

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/gracew/widget/graph/model"
	"github.com/stretchr/testify/assert"
)

func TestGenerateCode(t *testing.T) {
	api := model.API{Definition: &model.APIDefinition{
		Fields: []*model.FieldDefinition{
			&model.FieldDefinition{Name: "foo", Type: model.TypeBoolean},
			&model.FieldDefinition{Name: "bar", Type: model.TypeFloat},
			&model.FieldDefinition{Name: "baz", Type: model.TypeInt},
			&model.FieldDefinition{Name: "qux", Type: model.TypeString},
			&model.FieldDefinition{Name: "camelCase", Type: model.TypeString},
		},
	}}

	l := Launcher{API: api}
	generated, err := l.generateCode()
	assert.NoError(t, err)

	f, err := os.Open("test/generated.go")
	assert.NoError(t, err)

	expected, err := ioutil.ReadAll(f)
	assert.NoError(t, err)

	assert.Equal(t, string(expected), generated)
}
