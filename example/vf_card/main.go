package main

import (
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

	reciver := lib.NewVerifyCardReceiver(b)
	_, err = reciver.Rejected(func(event *lib.EventTransmitter) error {
		log.Println("vf card got rejected event :",event.ID())
		err := event.Approval()
		if err != nil{
			log.Println(err)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	_, err = reciver.Pending(func(event *lib.EventTransmitter) error {
		log.Println("vf_card got pending event :",event.ID())
		if event.ID() == "3" {
			err = event.Rejected()
			if err != nil {
				log.Println(err)
			}
			return nil
		}
		err = event.Approval()
		if err != nil {
			log.Println(err)
		}

		return nil
	})
	if err != nil {
		panic(err)
	}

	for i := 0; i < 6; i++ {
		time.Sleep(30 * time.Second)
	}

	log.Println("good jober")
}
