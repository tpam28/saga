package main

import (
	"flag"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/tpam28/saga/domain"
	"github.com/tpam28/saga/factory"
	"github.com/tpam28/saga/helper"
	"github.com/tpam28/saga/parser"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

func main() {
	flag.String("path", "1234", "help message for flagname")
	flag.String("output", "example.go", "help message for flagname")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	path := viper.GetString("path")
	output := viper.GetString("output")
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("mustn't read file:", err)
	}
	f := make(map[domain.Root][]map[string]map[string]interface{})

	err = yaml.Unmarshal(data, f)
	if err != nil {
		log.Fatal(err)
	}

	packageName := helper.FindPackageName(output)

	e, err := parser.ParseConfigSlice(f[domain.Milestone])
	if err != nil {
		log.Fatal(err)
	}
	fa := &factory.Factory{PathFile: output, PackageName: packageName}
	fa.Do(e)
}
