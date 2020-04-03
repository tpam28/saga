package factory

import (
	"bytes"
	"go/format"
	"html/template"
	"io/ioutil"
	"os"

	"github.com/Masterminds/sprig"
	"github.com/tpam28/saga/domain"
)

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

	err = os.Remove(f.PathFile)
	if err != nil {
		panic("remove file failed with err: " + err.Error())
	}

	b, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(f.PathFile, b, 0755)
	if err != nil {
		panic(err)
	}

}
