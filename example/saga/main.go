package main

import (
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-plugins/broker/rabbitmq/v2"
	"github.com/tpam28/saga/example/saga/lib"
	"log"
	"time"
)

//run simple orchestrator
func main() {
	b:= rabbitmq.NewBroker(func(options *broker.Options) {
		options.Addrs=append(options.Addrs,"amqp://evgen:wZCfo9@127.0.0.1:5672/test")
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
