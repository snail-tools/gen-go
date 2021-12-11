package gengo

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/snail-tools/gen-go/namer"
)

func TestDump(t *testing.T) {
	nr := namer.NewRawNamer("", namer.NewDefaultImportTracker())
	dump := NewDumper(nr)

	type DumpStruct struct {
		Name     string `db:"f_name"`
		Password string `db:"f_password"`
		isDelete bool   `db:"f_delete"`
		Age      int    `db:"f_age"`
	}
	t.Run("TypeOf", func(t *testing.T) {
		a := DumpStruct{
			Name:     "snail",
			Password: "abcdef",
			isDelete: true,
			Age:      10,
		}
		fmt.Println(dump.TypeOf(reflect.ValueOf(a).Type()))
	})
	t.Run("ValueOf", func(t *testing.T) {
		fmt.Println(dump.ValueOf("S"))
		fmt.Println(dump.ValueOf(123))
		fmt.Println(dump.ValueOf(123.111))
		fmt.Println(dump.ValueOf([]string{"a", "b", "c"}))
		fmt.Println(dump.ValueOf(map[string][]string{
			"aaa": {"a", "b", "c"},
			"bbb": {"a", "b", "c"},
		}))
		fmt.Println(dump.ValueOf(DumpStruct{
			Name:     "snail",
			Password: "abcdef",
			isDelete: true,
			Age:      10,
		}))
	})

}
