package types

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"
)

type TypeName interface {
	Pkg() *types.Package
	Name() string
	String() string
	Exported() bool
}

var _ TypeName = &types.TypeName{}

func Ref(pkgPath string, name string) TypeName {
	return &ref{pkgPath: pkgPath, name: name}
}

func ParseRef(ref string) (TypeName, error) {
	parts := strings.Split(ref, ".")
	if len(parts) == 1 {
		return nil, fmt.Errorf("unsupported ref: %s", ref)
	}
	return Ref(strings.Join(parts[0:len(parts)-1], "."), parts[len(parts)-1]), nil
}

func MustParseRef(ref string) TypeName {
	r, err := ParseRef(ref)
	if err != nil {
		panic(nil)
	}
	return r
}

type ref struct {
	pkgPath string
	name    string
}

func (ref) Underlying() types.Type {
	return nil
}

func (r *ref) String() string {
	return r.pkgPath + "." + r.name
}

func (r *ref) Pkg() *types.Package {
	return types.NewPackage(r.pkgPath, "")
}

func (r *ref) Name() string {
	return r.name
}

func (r *ref) Exported() bool {
	return ast.IsExported(r.name)
}
