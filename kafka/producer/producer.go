package producer

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func withDefaults() *kafka.ConfigMap {
	return &kafka.ConfigMap{"bootstrap.servers": "localhost"}
}

func Produce(topicPartition *kafka.TopicPartition, message string) {
	p, err := kafka.NewProducer(withDefaults())
	if err != nil {
		log.Print("An error occured while creating new producer")
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

	p.Produce(&kafka.Message{
		TopicPartition: *topicPartition,
		Value:          []byte(message),
	}, nil)

	// Wait for message deliveries before shutting down
	p.Flush(15 * 1000)
}
