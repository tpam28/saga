package main

import (
	"fmt"
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-plugins/broker/rabbitmq/v2"
	"github.com/tpam28/saga/example/saga/lib"
	"log"
	"time"
)

func main() {
	b:= rabbitmq.NewBroker(func(options *broker.Options) {
		options.Addrs=append(options.Addrs,"amqp://evgen:wZCfo9@127.0.0.1:5672/test")
	})

	err := b.Connect()
	if err != nil {
		panic(err)
	}

	transmitter := lib.NewVerifyConsumerTransmitter(b)
	time.Sleep(5*time.Second)
	for i := 0; i < 10; i++ {
		err = transmitter.Approval(lib.NewMessage(fmt.Sprint(i)))
		if err != nil {
			log.Println(err)
		}
		fmt.Println("vf consumer approval:",fmt.Sprint(i))
		time.Sleep(1*time.Second)
	}
	log.Println("good jober")
}
