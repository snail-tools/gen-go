package gengo

import (
	"fmt"
	"testing"

	gengotypes "github.com/snail-tools/gen-go/types"
	"github.com/stretchr/testify/require"
)

type modelGen struct {
	SnippetWriter
}

func (modelGen) New() Generator {
	return &modelGen{}
}

func (g *modelGen) Init(sw SnippetWriter) {
	g.SnippetWriter = sw
}

func (g *modelGen) Generate() {
	g.generateTable()
	g.generateFieldKeys()
	g.generateIndex()
}

func (g *modelGen) generateTable() {
	g.Do(`
func ([[ .typeName ]]) TableName() string {
	return "t_"+[[ "github.com/snail-tools/strcase.SnakeCase" | pkg ]]("[[ .typeName ]]")
}`, Args{
		"typeName": "User",
	})
}
func (g *modelGen) generateFieldKeys() {
	typeName := "User"
	list := []string{"UserID", "Name", "Email", "Password"}
	for _, n := range list {
		g.Do(`
func([[ .typeName ]]) FieldKey[[ .fieldName ]]() string {
    return [[ .fieldName | quote ]]
}
		`, Args{
			"typeName":  typeName,
			"fieldName": n,
		})
	}
}

func (g *modelGen) generateIndex() {
	indexs := map[string][]string{
		"i_org_id": {
			"OrgID",
		},
		"i_user_id": {
			"UserID",
		},
	}

	g.Do(`
[[ if .hasPrimary ]] func([[ .typeName ]]) Primary() []string {
	return [[ .primary ]]
} [[ end ]]

[[ if .hasIndexes ]] func([[ .typeName ]]) Indexes() [[ "github.com/snail-tools/sqlx.Indexes" | pkg ]] {
    return [[ .indexes ]]
} [[ end ]]
`, Args{
		"hasPrimary": true,
		"hasIndexes": true,
		"typeName":   "User",
		"primary":    g.Dumper().ValueOf([]string{"ID"}),
		"indexes":    g.Dumper().ValueOf(indexs),
	})
}

func TestGenfile(t *testing.T) {
	pwd := "github.com/snail-tools/gen-go/testdata/models"
	u, _, _ := gengotypes.Load(pwd)
	pkg := u.Package(pwd)
	file := NewGenFile(pkg)
	file.Generator(&modelGen{})
	data, err := file.Source()
	require.NoError(t, err)
	fmt.Println(string(data))
}
