package main

import (
	"flag"
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-plugins/broker/rabbitmq/v2"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/tpam28/saga/example/saga/lib"
	"log"
	"strings"
	"time"
)

func main() {
	config := viper.New()
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	config.AutomaticEnv()
	flag.Int("task", 1, "the number of task")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	b:= rabbitmq.NewBroker(func(options *broker.Options) {
		options.Addrs=append(options.Addrs,	config.GetString("rabbitmq.dsn"))
	})
	err := b.Connect()
	if err != nil {
		panic(err)
	}

	reciver := lib.NewConfirmTicketReceiver(b)
	_, err = reciver.Rejected(func(event *lib.EventTransmitter) error {
		if config.GetInt("task") == 3{
			if event.Retry() < 3{
				log.Println("confirm ticket got reject message event id",event.ID(),"event retry:",event.Retry())
				err = event.Reject()
				if err != nil{
					log.Println(err)
				}
				return nil
			}
		}

		err = event.Approve()
		if err != nil{
			log.Println(err)
		}

		log.Println("confirm_ticket got rejected event :", event.ID())

		return nil
	})
	if err != nil {
		panic(err)
	}
	_, err = reciver.Pending(func(event *lib.EventTransmitter) error {
		log.Println("confirm_ticket got pending event :", event.ID())
		if event.ID() == "3"{
			err = event.Reject()
			if err != nil {
				log.Println(err)
			}
			return nil
		}
		err = event.Approve()
		if err != nil {
			log.Println(err)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	time.Sleep(30 * time.Second)
}
