package gengo

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/scanner"
	"go/token"
	"io"
	"sort"
	"strings"

	"github.com/snail-tools/gen-go/namer"
	gengotypes "github.com/snail-tools/gen-go/types"
	"golang.org/x/tools/imports"

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

func (f *GenFile) Source() ([]byte, error) {
	buf := &bytes.Buffer{}
	buf.WriteString(`package ` + strcase.SnakeCase(f.pkg.Pkg().Name()) + `
`)

	writeImports(buf, f.imports.Imports())

	if _, err := io.Copy(buf, f.body); err != nil {
		return nil, err
	}

	data := buf.Bytes()
	lines := bytes.Split(data, []byte("\n"))

	if _, err := parser.ParseFile(token.NewFileSet(), "", data, parser.AllErrors); err != nil {
		if sl, ok := err.(scanner.ErrorList); ok {
			for i := range sl {
				l := sl[i].Pos.Line

				fmt.Println(sl[i].Pos)

				for i := l - 3; i < l; i++ {
					if i > 0 {
						fmt.Printf("%d\t%s\n", i+1, string(lines[i]))
					}
				}

				col := sl[i].Pos.Column - 1
				if col < 0 {
					col = 0
				}
				fmt.Printf("\t%s↑\n", strings.Repeat(" ", col))
				fmt.Println(sl[i].Msg)
			}
		}
		return nil, err
	}

	return imports.Process("", data, &imports.Options{
		FormatOnly: true,
	})
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
