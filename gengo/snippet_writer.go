package gengo

import (
	"bytes"
	"fmt"
	"io"
	"runtime"
	"strconv"
	"text/template"

	"github.com/snail-tools/gen-go/namer"
	gengotypes "github.com/snail-tools/gen-go/types"
)

func Snippet(format string, args ...Args) func(s SnippetWriter) {
	return func(s SnippetWriter) {
		s.Do(format, args...)
	}
}

type SnippetWriter interface {
	io.Writer
	Do(format string, args ...Args)
	Dumper() *Dumper
}

type Args = map[string]interface{}

func NewSnippetWriter(w io.Writer, ns namer.NameSystems) SnippetWriter {
	sw := &snippetWriter{
		Writer: w,
		ns:     ns,
	}
	return sw
}

func createRender(ns namer.NameSystems) func(r func(s SnippetWriter)) string {
	return func(r func(s SnippetWriter)) string {
		b := bytes.NewBuffer(nil)
		r(NewSnippetWriter(b, ns))
		return b.String()
	}
}

type snippetWriter struct {
	io.Writer
	ns namer.NameSystems
}

func (s *snippetWriter) Do(format string, args ...Args) {
	_, file, line, _ := runtime.Caller(1)

	tmpl, err := template.
		New(fmt.Sprintf("%s:%d", file, line)).
		Delims("[[", "]]").
		Funcs(createFuncMap(s.ns)).
		Parse(format)

	if err != nil {
		panic(err)
	}

	finalArgs := Args{}

	for i := range args {
		a := args[i]
		for k := range a {
			finalArgs[k] = a[k]
		}
	}

	if err := tmpl.Execute(s.Writer, finalArgs); err != nil {
		panic(err)
	}
}

func (s *snippetWriter) Dumper() *Dumper {
	if rawNamer, ok := s.ns["raw"]; ok {
		return NewDumper(rawNamer)
	}
	return nil
}

func createFuncMap(nameSystems namer.NameSystems) template.FuncMap {
	funcMap := template.FuncMap{}

	funcMap["pkg"] = createPackage(nameSystems)
	funcMap["render"] = createRender(nameSystems)
	funcMap["quote"] = strconv.Quote

	return funcMap
}

func createPackage(nameSystems namer.NameSystems) func(v interface{}) string {
	return func(v interface{}) string {
		switch x := v.(type) {
		case string:
			ref, err := gengotypes.ParseRef(x)
			if err != nil {
				return x
			}
			return nameSystems["raw"].Name(ref)
		case gengotypes.TypeName:
			return nameSystems["raw"].Name(x)
		default:
			panic("unspported")
		}
	}
}
