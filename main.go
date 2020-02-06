package main

import (
	"flag"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/tpam28/saga/domain"
	"github.com/tpam28/saga/factory"
	"github.com/tpam28/saga/parser"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

func main() {
	// using standard library "flag" package
	flag.String("path", "1234", "help message for flagname")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	path := viper.GetString("path") // retrieve value from viper
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("mustn't read file:", err)
	}
	f := make(map[domain.Root][]map[string]map[string]interface{})

	err = yaml.Unmarshal(data, f)
	if err != nil {
		log.Fatal(err)
	}


	e,err :=parser.ParseConfigSlice(f[domain.Milestone])
	if err != nil{
		log.Fatal(err)
	}
	fa:= &factory.Factory{}
	fa.Do(e)
}

