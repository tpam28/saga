package factory

import (
	"bytes"
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
	File *os.File
}

func (f *Factory) Write() error {
	_, err := f.File.Write([]byte(imports))
	if err != nil {
		return err
	}
	return nil
}

func (f *Factory) Do(list domain.StepList) {
	t := template.New("const.tpl")
	t = t.Funcs(sprig.FuncMap())
	t, err := t.ParseFiles("factory/const.tpl")
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, list)
	if err != nil {
		panic(err)
	}
	err =ioutil.WriteFile("example.go",buf.Bytes(),0755)
	if err != nil {
		panic(err)
	}

}
