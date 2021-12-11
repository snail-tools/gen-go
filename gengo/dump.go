package gengo

import (
	"bytes"
	"fmt"
	"go/ast"
	"reflect"
	"sort"
	"strconv"

	"github.com/snail-tools/gen-go/namer"
	gengotypes "github.com/snail-tools/gen-go/types"
)

func NewDumper(rawNamer namer.Namer) *Dumper {
	return &Dumper{
		namer: rawNamer,
	}
}

type Dumper struct {
	namer namer.Namer
}

func (d *Dumper) Name(named gengotypes.TypeName) string {
	return d.namer.Name(named)
}

func (d *Dumper) TypeOf(tpe reflect.Type) string {
	if tpe.PkgPath() != "" {
		return d.Name(gengotypes.Ref(tpe.PkgPath(), tpe.Name()))
	}

	switch tpe.Kind() {
	case reflect.Ptr:
		return "*" + d.TypeOf(tpe.Elem())
	case reflect.Chan:
		return "chan " + d.TypeOf(tpe.Elem())
	case reflect.Struct:
		b := bytes.NewBufferString("struct {")

		for i := 0; i < tpe.NumField(); i++ {
			f := tpe.Field(i)

			if !f.Anonymous {
				_, _ = fmt.Fprintf(b, "%s ", f.Name)
			}

			b.WriteString(d.TypeOf(f.Type))

			if tag := f.Tag; tag != "" {
				_, _ = fmt.Fprintf(b, " `%s`", tag)
			}

			b.WriteString("\n")
		}

		b.WriteString("}")

		return b.String()
	case reflect.Array:
		return fmt.Sprintf("[%d]%s", tpe.Len(), d.TypeOf(tpe.Elem()))
	case reflect.Slice:
		return fmt.Sprintf("[]%s", d.TypeOf(tpe.Elem()))
	case reflect.Map:
		return fmt.Sprintf("map[%s]%s", d.TypeOf(tpe.Key()), d.TypeOf(tpe.Elem()))
	default:
		return tpe.String()
	}
}

var basicKinds = map[reflect.Kind]bool{
	reflect.Bool:       true,
	reflect.Int:        true,
	reflect.Int8:       true,
	reflect.Int16:      true,
	reflect.Int32:      true,
	reflect.Int64:      true,
	reflect.Uint:       true,
	reflect.Uint8:      true,
	reflect.Uint16:     true,
	reflect.Uint32:     true,
	reflect.Uint64:     true,
	reflect.Uintptr:    true,
	reflect.Float32:    true,
	reflect.Float64:    true,
	reflect.Complex64:  true,
	reflect.Complex128: true,
}

func (d *Dumper) ValueOf(in interface{}) string {
	rv, ok := in.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(in)
	}

	if rv.Kind() == reflect.Ptr && rv.IsNil() {
		return "nil"
	}
	tpe := rv.Type()
	switch tpe.Kind() {
	case reflect.Ptr:
		kind := rv.Elem().Kind()
		if _, ok := basicKinds[kind]; ok {
			return fmt.Sprintf("func(v %s) *%s { return &v }(%s)", kind, kind, d.ValueOf(rv.Elem()))
		}
		return fmt.Sprintf("&(%s)", d.ValueOf(rv.Elem()))
	case reflect.Struct:
		buf := bytes.NewBufferString(d.ReflectTypeOf(tpe))
		buf.WriteString(`{`)

		c := 0

		for i := 0; i < rv.NumField(); i++ {
			f := rv.Field(i)
			ft := tpe.Field(i)
			if ast.IsExported(ft.Name) {
				v := d.ValueOf(f)

				if v == "" {
					continue
				}

				if c == 0 {
					buf.WriteString("\n")
				}

				buf.WriteString(ft.Name)
				buf.WriteString(":")
				buf.WriteString(v)
				buf.WriteString(",")
				buf.WriteString("\n")

				c++
			}
		}

		buf.WriteString(`}`)

		return buf.String()
	case reflect.Map:
		buf := bytes.NewBufferString(d.ReflectTypeOf(tpe))
		buf.WriteString(`{`)

		keyLits := make([]string, 0)
		keyValues := map[string]reflect.Value{}

		for _, key := range rv.MapKeys() {
			k := d.ValueOf(key)
			keyLits = append(keyLits, k)
			keyValues[k] = rv.MapIndex(key)
		}

		sort.Strings(keyLits)

		for i, k := range keyLits {
			if i == 0 {
				buf.WriteString("\n")
			}

			buf.WriteString(k)
			buf.WriteString(":")
			buf.WriteString(d.ValueOf(keyValues[k]))
			buf.WriteString(",")
			buf.WriteString("\n")
		}

		buf.WriteString(`}`)
		return buf.String()
	case reflect.Slice, reflect.Array:
		buf := bytes.NewBufferString(d.ReflectTypeOf(tpe))
		buf.WriteString(`{`)

		for i := 0; i < rv.Len(); i++ {
			if i == 0 {
				buf.WriteString("\n")
			}

			buf.WriteString(d.ValueOf(rv.Index(i)))
			buf.WriteString(",")
			buf.WriteString("\n")
		}

		buf.WriteString(`}`)

		return buf.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int64:
		return fmt.Sprintf("%d", rv.Int())
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
		return fmt.Sprintf("%d", rv.Uint())
	case reflect.Int32:
		if b, ok := rv.Interface().(rune); ok {
			r := strconv.QuoteRune(b)
			if len(r) == 3 {
				return r
			}
		}
		return fmt.Sprintf("%d", rv.Int())
	case reflect.Bool:
		return strconv.FormatBool(rv.Bool())
	case reflect.Float32:
		return strconv.FormatFloat(rv.Float(), 'f', -1, 32)
	case reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'f', -1, 64)
	case reflect.String:
		return strconv.Quote(rv.String())
	// case reflect.Interface:
	case reflect.Invalid:
		return "nil"
	default:
		panic(fmt.Errorf("%s is an unsupported type", tpe.String()))
	}
}

func (d *Dumper) ReflectTypeOf(tpe reflect.Type) string {
	return d.TypeOf(tpe)
}
