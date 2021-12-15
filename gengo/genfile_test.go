package gengo

import (
	"fmt"
	"path/filepath"
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
		"typeName": "Account",
	})
}
func (g *modelGen) generateFieldKeys() {
	typeName := "Account"
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

[[ if .hasIndexes ]] func([[ .typeName ]]) Indexes() map[string][]string {
    return [[ .indexes ]]
} [[ end ]]
`, Args{
		"hasPrimary": true,
		"hasIndexes": true,
		"typeName":   "Account",
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
	filename := filepath.Dir(pkg.GoFiles()[0]) + "/account__generated.go"
	fmt.Println(filename)
	err := file.WriteFile(filename)
	require.NoError(t, err)
}
