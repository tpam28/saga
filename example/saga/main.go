package main

import (
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-plugins/broker/rabbitmq/v2"
	"github.com/spf13/viper"
	"github.com/tpam28/saga/example/saga/lib"
	"log"
	"strings"
	"time"
)

//run simple orchestrator
func main() {
	config := viper.New()
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	config.AutomaticEnv()

	b:= rabbitmq.NewBroker(func(options *broker.Options) {
		options.Addrs=append(options.Addrs, config.GetString("rabbitmq.dsn"))
	})
	err := b.Connect()
	if err !=nil{
		panic(err)
	}
	log.Println(b.Address())
	orchestrator := lib.NewOrchestrator(b, nil)
	_, err = orchestrator.Do()
	if err != nil {
		panic(err)
	}
	log.Println("good job")
	for i:=0; i<5;i++{
		time.Sleep(30 * time.Second)
	}
}
