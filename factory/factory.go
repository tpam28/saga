package factory

import (
	"bytes"
	"go/format"
	"html/template"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/Masterminds/sprig"
	"github.com/tpam28/saga/domain"
)

var imports = `
import (
	"github.com/micro/go-micro/v2/broker"
)
`

var fns = template.FuncMap{
	"last": func(x int, a interface{}) bool {
		return x == reflect.ValueOf(a).Len()-1
	},
}

type Factory struct {
	domain.StepList
	PathFile    string
	PackageName string
}

func (f *Factory) Do(list domain.StepList) {
	t := template.New("const.tpl")
	t = t.Funcs(sprig.FuncMap())
	t, err := t.ParseFiles("factory/const.tpl")
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	buf.Write([]byte("//Automatically generated file; DO NOT EDIT\n"))
	buf.Write([]byte("package " + f.PackageName ))
	err = t.Execute(&buf, list)
	if err != nil {
		panic(err)
	}
	os.Remove(f.PathFile)
	b, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(f.PathFile, b, 0755)
	if err != nil {
		panic(err)
	}

}
