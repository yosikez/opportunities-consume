package main

import (
	"log"

	"github.com/yosikez/opportunities-consume/controller"
	"github.com/yosikez/opportunities-consume/database"
	"github.com/yosikez/opportunities-consume/rabbitmq"
)

func main() {
	err := database.Connect()

	if err != nil {
		log.Fatalf("failed to connect to database : %v", err)
	}

	rmqCfg, rmq, err := rabbitmq.NewRabbitMQ()

	if err != nil {
		log.Fatalf("failed to connect to rabbitmq : %v", err)
	}

	defer rmq.Connection.Close()
	defer rmq.Channel.Close()

	err = rmq.Channel.ExchangeDeclare(
		rmqCfg.ExchangeName,
		rmqCfg.ExchangeKind,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatalf("failed to declare exchange : %v", err)
	}

	consumeController := controller.NewConsumeController(rmq, rmqCfg)


	forever := make(chan struct{})
	go func ()  {
		consumeController.StartConsumer()
	}()

	log.Printf("Waiting for incoming data from RabbitMQ....")
	<-forever
}
