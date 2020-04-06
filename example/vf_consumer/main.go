package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-plugins/broker/rabbitmq/v2"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/tpam28/saga/example/saga/lib"
)

func main() {
	config := viper.New()
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	config.AutomaticEnv()
	flag.Int("task", 1, "the number of task")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	b := rabbitmq.NewBroker(func(options *broker.Options) {
		options.Addrs = append(options.Addrs, config.GetString("rabbitmq.dsn"))
	})

	err := b.Connect()
	if err != nil {
		panic(err)
	}
	reciever := lib.NewVerifyConsumerReceiver(b)
	reciever.Rejected(func(event *lib.EventTransmitter) error {
		if config.GetInt("task") == 4 && event.Retry() < 3 {
			log.Println("vf consumer got rejected event", event.ID(), "attention:", event.Retry())

			err := event.Reject()
			if err != nil {
				log.Println(err)
			}
			return nil
		}
		log.Println("vf consumer got rejected event", event.ID())

		err := event.Approve()
		if err != nil {
			log.Println(err)
		}
		return nil
	})

	transmitter := lib.NewVerifyConsumerTransmitter(b)
	time.Sleep(5 * time.Second)
	err = transmitter.Approve(lib.NewMessage(fmt.Sprint(viper.GetString("task"))))
	fmt.Println("vf consumer approval:", fmt.Sprint(viper.GetString("task")))
	if err != nil {
		log.Println(err)
	}
	time.Sleep(30 * time.Second)

}
