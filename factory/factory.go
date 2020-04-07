package factory

import (
	"bytes"
	"github.com/Masterminds/sprig"
	"github.com/tpam28/saga/domain"
	"go/format"
	"html/template"
	"io/ioutil"
	"log"
	"os"
)

type Factory struct {
	domain.StepList
	PathFile    string
	PackageName string
}

var isNil = template.FuncMap{
	"isNil": func(i *domain.Step) bool { return i == nil },
}

func (f *Factory) Do(list domain.StepList) {
	t := template.New("const.tpl")
	t = t.Funcs(sprig.FuncMap())
	t = t.Funcs(isNil)
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
		log.Println("remove file failed with err: " + err.Error())
	}

	b, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(f.PathFile, b, 0755)
	if err != nil{
		panic(err)
	}

}
