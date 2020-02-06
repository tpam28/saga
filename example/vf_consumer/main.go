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
		err := event.Approval()
		if err != nil {
			log.Println(err)
		}
		return nil
	})

	transmitter := lib.NewVerifyConsumerTransmitter(b)
	time.Sleep(5 * time.Second)
	err = transmitter.Approval(lib.NewMessage(fmt.Sprint(viper.GetString("task"))))
	if err != nil {
		log.Println(err)
	}

	fmt.Println("vf consumer approval:", fmt.Sprint(viper.GetString("task")))
}
