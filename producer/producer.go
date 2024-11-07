package producer

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Producer interface {
	Produce(topicPartition *kafka.TopicPartition, message string)
}

type ConfluentIncKafkaProducer struct {
	cfg kafka.ConfigMap
}

func New(cfg kafka.ConfigMap) Producer {
	return &ConfluentIncKafkaProducer{cfg: cfg}
}

func (k *ConfluentIncKafkaProducer) Produce(topicPartition *kafka.TopicPartition, message string) {
	p, err := kafka.NewProducer(&k.cfg)
	if err != nil {
		log.Println("An error occurred while creating a new producer:", err)
		return
	}
	defer p.Close()

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					log.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	err = p.Produce(&kafka.Message{
		TopicPartition: *topicPartition,
		Value:          []byte(message),
	}, nil)
	if err != nil {
		log.Println("An error occurred producing message:", err)
		return
	}
	// Wait for message deliveries before shutting down
	p.Flush(15 * 1000)
}
