package kafka

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Grealish17/parvpo/internal/model"

	"github.com/IBM/sarama"
)

type ConsumerGroup struct {
	ready chan bool
}

func NewConsumerGroup() ConsumerGroup {
	return ConsumerGroup{
		ready: make(chan bool),
	}
}

func (consumer *ConsumerGroup) Ready() <-chan bool {
	return consumer.ready
}

func (consumer *ConsumerGroup) Setup(_ sarama.ConsumerGroupSession) error {
	close(consumer.ready)

	return nil
}

func (consumer *ConsumerGroup) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *ConsumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():

			pm := model.RequestMessage{}
			err := json.Unmarshal(message.Value, &pm)
			if err != nil {
				fmt.Println("Consumer group error", err)
			}

			log.Println(
				pm.ID,
				pm.UserEmail,
				pm.Price,
				pm.HomeTeam,
				pm.AwayTeam,
				pm.DateTime,
			)

			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}
