package types

import (
	"go/ast"
	"go/token"
	"go/types"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

type Package interface {
	Pkg() *types.Package
	Module() *packages.Module
	SourceDir() string
	Files() []*ast.File
	GoFiles() []string

	Doc(pos token.Pos) (map[string][]string, []string)
	Comment(pos token.Pos) []string

	// Type
	Type(name string) *types.TypeName
	Types() map[string]*types.TypeName
	// Constant
	Constant(name string) *types.Const
	Constants() map[string]*types.Const
	// Function
	Function(name string) *types.Func
	Functions() map[string]*types.Func
	// Methods
	MethodsOf(n *types.Named, canPtr bool) []*types.Func
	// ResultsOf
	ResultsOf(tpe *types.Func) (results Results, resultN int)
	// Position
	Position(pos token.Pos) token.Position

	Eval(expr ast.Expr) (types.TypeAndValue, error)
}

func newPkg(pkg *packages.Package, u Universe) Package {
	pi := &pkgInfo{
		u: u,

		pkg: pkg,

		endLineToCommentGroup:         map[fileLine]*ast.CommentGroup{},
		endLineToTrailingCommentGroup: map[fileLine]*ast.CommentGroup{},

		signatures:  map[*types.Signature]ast.Node{},
		funcResults: map[*types.Signature][]TypeAndValues{},

		constants: map[string]*types.Const{},
		types:     map[string]*types.TypeName{},
		funcs:     map[string]*types.Func{},

		methods: map[*types.Named][]*types.Func{},
	}

	fileLineFor := func(pos token.Pos, deltaLine int) fileLine {
		position := pi.pkg.Fset.Position(pos)
		return fileLine{position.Filename, position.Line + deltaLine}
	}

	collectCommentGroup := func(c *ast.CommentGroup, isTrailing bool, stmtPos token.Pos) {
		fl := fileLineFor(stmtPos, 0)

		if c != nil && c.Pos() == stmtPos {
			// stmt is CommentGroup
			fl = fileLineFor(c.End(), 0)
		} else {
			fl = fileLineFor(stmtPos, -1)
		}

		if isTrailing {
			if cc := pi.endLineToTrailingCommentGroup[fl]; cc == nil {
				pi.endLineToTrailingCommentGroup[fl] = c
			}
		} else {
			if cc := pi.endLineToCommentGroup[fl]; cc == nil {
				pi.endLineToCommentGroup[fl] = c
			}
		}
	}

	for i := range pi.pkg.Syntax {
		f := pi.pkg.Syntax[i]

		ast.Inspect(f, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.CallExpr:
				// signature will be from other package
				// stored p.TypesInfo.Uses[*ast.Ident].(*types.PkgName)
				fn := pi.pkg.TypesInfo.TypeOf(x.Fun)
				if fn != nil {
					if s, ok := fn.(*types.Signature); ok {
						if n, ok := pi.signatures[s]; ok {
							switch n.(type) {
							case *ast.FuncDecl, *ast.FuncLit:
								// skip declared functions
							default:
								pi.signatures[s] = x.Fun
							}
						} else {
							pi.signatures[s] = x.Fun
						}
					}
				}
			case *ast.FuncDecl:
				fn := pi.pkg.TypesInfo.TypeOf(x.Name)
				if fn != nil {
					pi.signatures[fn.(*types.Signature)] = x
				}
			case *ast.FuncLit:
				fn := pi.pkg.TypesInfo.TypeOf(x)
				if fn != nil {
					pi.signatures[fn.(*types.Signature)] = x
				}
			case *ast.CommentGroup:
				collectCommentGroup(x, false, x.Pos())
			case *ast.ValueSpec:
				collectCommentGroup(x.Doc, false, x.Pos())
				collectCommentGroup(x.Comment, true, x.Pos())
			case *ast.ImportSpec:
				collectCommentGroup(x.Doc, false, x.Pos())
				collectCommentGroup(x.Comment, true, x.Pos())
			case *ast.TypeSpec:
				collectCommentGroup(x.Doc, false, x.Pos())
				collectCommentGroup(x.Comment, true, x.Pos())
			case *ast.Field:
				collectCommentGroup(x.Doc, false, x.Pos())
				collectCommentGroup(x.Comment, true, x.Pos())
			}
			return true
		})
	}

	for ident := range pi.pkg.TypesInfo.Defs {
		switch x := pi.pkg.TypesInfo.Defs[ident].(type) {
		case *types.Func:
			s := x.Type().(*types.Signature)

			if r := s.Recv(); r != nil {
				var named *types.Named

				switch t := r.Type().(type) {
				case *types.Pointer:
					if n, ok := t.Elem().(*types.Named); ok {
						named = n
					}
				case *types.Named:
					named = t
				}

				if named != nil {
					pi.methods[named] = append(pi.methods[named], x)
				}
			} else {
				pi.funcs[x.Name()] = x
			}
		case *types.TypeName:
			pi.types[x.Name()] = x
		case *types.Const:
			pi.constants[x.Name()] = x
		}
	}

	return pi
}

type pkgInfo struct {
	u   Universe
	pkg *packages.Package

	constants map[string]*types.Const
	types     map[string]*types.TypeName
	funcs     map[string]*types.Func
	methods   map[*types.Named][]*types.Func

	endLineToCommentGroup         map[fileLine]*ast.CommentGroup
	endLineToTrailingCommentGroup map[fileLine]*ast.CommentGroup

	signatures  map[*types.Signature]ast.Node
	funcResults map[*types.Signature][]TypeAndValues
}

func (pi *pkgInfo) SourceDir() string {
	if pi.pkg.PkgPath == pi.Module().Path {
		return pi.Module().Dir
	}
	return filepath.Join(pi.Module().Dir, pi.pkg.PkgPath[len(pi.Module().Path):])
}

func (pi *pkgInfo) Pkg() *types.Package {
	return pi.pkg.Types
}

func (pi *pkgInfo) Module() *packages.Module {
	return pi.pkg.Module
}

func (pi *pkgInfo) Files() []*ast.File {
	return pi.pkg.Syntax
}

func (pi *pkgInfo) Type(name string) *types.TypeName {
	return pi.types[name]
}

func (pi *pkgInfo) Types() map[string]*types.TypeName {
	return pi.types
}

func (pi *pkgInfo) Constant(name string) *types.Const {
	return pi.constants[name]
}

func (pi *pkgInfo) Constants() map[string]*types.Const {
	return pi.constants
}

func (pi *pkgInfo) Function(n string) *types.Func {
	return pi.funcs[n]
}

func (pi *pkgInfo) Functions() map[string]*types.Func {
	return pi.funcs
}

func (pi *pkgInfo) MethodsOf(n *types.Named, ptr bool) []*types.Func {
	funcs, _ := pi.methods[n]

	if ptr {
		return funcs
	}

	notPtrMethods := make([]*types.Func, 0)

	for i := range funcs {
		s := funcs[i].Type().(*types.Signature)

		if _, ok := s.Recv().Type().(*types.Pointer); !ok {
			notPtrMethods = append(notPtrMethods, funcs[i])
		}
	}

	return notPtrMethods
}

func (pi *pkgInfo) Position(pos token.Pos) token.Position {
	return pi.pkg.Fset.Position(pos)
}

func (pi *pkgInfo) Eval(expr ast.Expr) (types.TypeAndValue, error) {
	return types.Eval(pi.pkg.Fset, pi.pkg.Types, expr.Pos(), StringifyNode(pi.pkg.Fset, expr))
}

func (pi *pkgInfo) Doc(pos token.Pos) (map[string][]string, []string) {
	return ExtractCommentTags(commentLinesFrom(pi.priorCommentLines(pos, -1)))
}

func (pi *pkgInfo) Comment(pos token.Pos) []string {
	return commentLinesFrom(pi.priorCommentLines(pos, 0))
}

func (pi *pkgInfo) GoFiles() []string {
	return pi.pkg.GoFiles
}

func (pi *pkgInfo) priorCommentLines(pos token.Pos, deltaLines int) *ast.CommentGroup {
	position := pi.pkg.Fset.Position(pos)
	key := fileLine{position.Filename, position.Line + deltaLines}
	if deltaLines == 0 {
		// should ignore trailing comments
		// when deltaLines eq 0 means find trailing comments
		if _, ok := pi.endLineToTrailingCommentGroup[key]; ok {
			return nil
		}
	}
	return pi.endLineToCommentGroup[key]
}

type fileLine struct {
	file string
	line int
}

func commentLinesFrom(commentGroups ...*ast.CommentGroup) (comments []string) {
	if len(commentGroups) == 0 {
		return nil
	}

	for _, commentGroup := range commentGroups {
		if commentGroup == nil {
			continue
		}

		for _, line := range strings.Split(strings.TrimSpace(commentGroup.Text()), "\n") {
			// 跳过 go: prefix
			if strings.HasPrefix(line, "go:") {
				continue
			}
			comments = append(comments, line)
		}
	}
	return comments
}
