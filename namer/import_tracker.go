package namer

import (
	"go/token"
	"strings"

	gengotypes "github.com/snail-tools/gen-go/types"
)

type ImportTracker interface {
	AddType(o gengotypes.TypeName)
	LocalNameOf(packagePath string) string
	PathOf(localName string) (string, bool)
	Imports() map[string]string
}

func NewDefaultImportTracker() ImportTracker {
	return &defaultImportTracker{
		pathToName: map[string]string{},
		nameToPath: map[string]string{},
	}
}

type defaultImportTracker struct {
	pathToName map[string]string
	nameToPath map[string]string
}

func (tracker *defaultImportTracker) AddType(o gengotypes.TypeName) {
	path := o.Pkg().Path()

	if _, ok := tracker.pathToName[path]; ok {
		return
	}

	localName := golangTrackerLocalName(path)

	tracker.nameToPath[localName] = path
	tracker.pathToName[path] = localName

}
func golangTrackerLocalName(name string) string {
	name = strings.Replace(name, "/", "_", -1)
	name = strings.Replace(name, ".", "_", -1)
	name = strings.Replace(name, "-", "_", -1)
	// 关键字
	if token.Lookup(name).IsKeyword() {
		name = "_" + name
	}
	return name
}

func (tracker *defaultImportTracker) LocalNameOf(packagePath string) string {
	return tracker.pathToName[packagePath]
}

func (tracker *defaultImportTracker) PathOf(localName string) (string, bool) {
	name, ok := tracker.nameToPath[localName]
	return name, ok
}

func (tracker *defaultImportTracker) Imports() map[string]string {
	return tracker.pathToName
}
