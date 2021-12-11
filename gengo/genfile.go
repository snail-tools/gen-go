package gengo

import (
	"bytes"
	"fmt"
	"io"
	"sort"

	"github.com/snail-tools/gen-go/namer"
	gengotypes "github.com/snail-tools/gen-go/types"

	"github.com/snail-tools/strcase"
)

type GenFile struct {
	pkg     gengotypes.Package
	body    *bytes.Buffer
	imports namer.ImportTracker
	sw      SnippetWriter
}

func NewGenFile(pkg gengotypes.Package) *GenFile {
	im := namer.NewDefaultImportTracker()
	return &GenFile{
		pkg:     pkg,
		body:    bytes.NewBuffer(nil),
		imports: im,
	}
}

func (g *GenFile) Generator(gen Generator) {
	g.sw = NewSnippetWriter(g.body, map[string]namer.Namer{
		"raw": namer.NewRawNamer(g.pkg.Pkg().Path(), g.imports),
	})

	gen.Init(g.sw)
	gen.Generate()
}

func (f *GenFile) Bytes() ([]byte, error) {
	buf := &bytes.Buffer{}
	buf.WriteString(`package ` + strcase.SnakeCase(f.pkg.Pkg().Name()) + `
`)

	writeImports(buf, f.imports.Imports())

	if _, err := io.Copy(buf, f.body); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func writeImports(w io.Writer, pathToName map[string]string) {
	importPaths := make([]string, 0)
	for p := range pathToName {
		importPaths = append(importPaths, p)
	}
	sort.Strings(importPaths)

	if len(importPaths) > 0 {
		_, _ = fmt.Fprintf(w, `
import (
`)

		for _, p := range importPaths {
			_, _ = fmt.Fprintf(w, `	%s "%s"
`, pathToName[p], p)
		}

		_, _ = fmt.Fprintf(w, `)
`)
	}
}
