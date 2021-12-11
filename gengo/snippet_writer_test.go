package gengo

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/snail-tools/gen-go/namer"
)

func TestSnippet(t *testing.T) {
	body := bytes.NewBuffer(nil)
	ims := namer.NewDefaultImportTracker()
	sw := NewSnippetWriter(body, map[string]namer.Namer{
		"raw": namer.NewRawNamer("github.com/snail-tools/gen-go/testdata/models", ims),
	})

	sw.Do(`
func([[ .typeName ]]) FieldKey[[ .fieldName ]]() string {
    return [[ .fieldName ]]
}
	`, Args{
		"typeName":  "Model",
		"fieldName": "FieldName",
	})

	fmt.Println("--------------------------------")
	fmt.Println(body.String())
	fmt.Println("--------------------------------")

}
