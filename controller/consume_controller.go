package controller

import (
	"encoding/json"
	"log"
	"time"

	"github.com/yosikez/opportunities-consume/config"
	"github.com/yosikez/opportunities-consume/database"
	"github.com/yosikez/opportunities-consume/input"
	"github.com/yosikez/opportunities-consume/model"
)

type ConsumeController struct {
	rmq    *config.RabbitMQConnection
	rmqCfg *config.RabbitMQ
}

func NewConsumeController(rmq *config.RabbitMQConnection, rmqCfg *config.RabbitMQ) *ConsumeController {
	return &ConsumeController{
		rmq:    rmq,
		rmqCfg: rmqCfg,
	}
}

func (cs *ConsumeController) StartConsumer() {
	cs.rmqDeclare("opportunity_create_queue")
}

func (cs *ConsumeController) rmqDeclare(queueName string) {
	q, err := cs.rmq.Channel.QueueDeclare(queueName, false, false, false, false, nil)

	if err != nil {
		log.Fatalf("failed to declare queue : %v", err)
	}

	err = cs.rmq.Channel.QueueBind(q.Name, queueName, cs.rmqCfg.ExchangeName, false, nil)

	if err != nil {
		log.Fatalf("failed to bind queue : %v", err)
	}

	msgs, err := cs.rmq.Channel.Consume(queueName, "", true, false, false, false, nil)

	if err != nil {
		log.Fatalf("failed to register a consumer : %v", err)
	}

	var input input.Opportunity

	for d := range msgs {
		err := json.Unmarshal(d.Body, &input)
		if err != nil {
			log.Fatalf("failed to unmarshal : %v", err)
		}
		cs.InsertData(&input)
	}
}

func (cs *ConsumeController) InsertData(data *input.Opportunity) {
	lastModified, err := time.Parse("2006-01-02 15:04:05", data.LastModified)
	if err != nil {
		log.Print("error parse time")
	}

	opportunity := &model.Opportunity{
		Code:            data.Code,
		ClientCode:      data.ClientCode,
		PicEmail:        data.PicEmail,
		OpportunityName: data.OpportunityName,
		Description:     data.Description,
		SalesEmail:      data.SalesEmail,
		Status:          data.Status,
		LastModified:    lastModified,
	}

	if err := database.DB.Create(&opportunity).Error; err != nil {
		log.Fatalf("failed to insert data opportunity : %v", err)
	}

	for _, resource := range data.Resources {

		r := &model.Resource{
			OpportunityId:   opportunity.Id,
			Qty:             resource.Qty,
			Position:        resource.Position,
			Level:           resource.Level,
			Ctc:             resource.Ctc,
			ProjectDuration: resource.ProjectDuration,
		}
		if err := database.DB.Create(&r).Error; err != nil {
			log.Fatalf("failed to insert data resource : %v", err)
		}
	}

	log.Println("success insert data")
}
