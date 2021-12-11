package namer

import "github.com/snail-tools/gen-go/types"

type Namer interface {
	Name(types.TypeName) string
}

type NameSystems map[string]Namer

type Names map[types.TypeName]string

func NewRawNamer(pkgPath string, tracker ImportTracker) Namer {
	return &rawNamer{pkgPath: pkgPath, tracker: tracker}
}

type rawNamer struct {
	pkgPath string
	tracker ImportTracker
	Names
}

func (n *rawNamer) Name(typeName types.TypeName) string {
	if n.Names == nil {
		n.Names = Names{}
	}

	if name, ok := n.Names[typeName]; ok {
		return name
	}

	pkgPath := typeName.Pkg().Path()

	if pkgPath == n.pkgPath {
		name := typeName.Name()
		if name != "" {
			return name
		}
		return typeName.String()
	} else {
		n.tracker.AddType(typeName)
		return n.tracker.LocalNameOf(pkgPath) + "." + typeName.Name()
	}
}
