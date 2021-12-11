package types

import (
	"fmt"
	"go/types"
	"testing"
)

func TestLoadModel(t *testing.T) {
	pwd := "github.com/snail-tools/gen-go/testdata/models"
	u, _, _ := Load(pwd)
	pkg := u.Package(pwd)
	t.Run("Doc", func(t *testing.T) {
		st := pkg.Type("Account")
		a1, a2 := pkg.Doc(st.Pos())
		fmt.Println(a2)

		aa := pkg.Comment(st.Pos())
		fmt.Println(aa)
		fmt.Println("name:", "Account")
		for _, s := range a1["def"] {
			fmt.Println(s)
		}

		structType := st.Type().Underlying().(*types.Struct)
		for i := 0; i < structType.NumFields(); i++ {
			field := structType.Field(i)
			fmt.Println(field.Name(), field.Type())
		}

	})
}
