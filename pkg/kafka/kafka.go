package kafka

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/IBM/sarama"
)

var producer sarama.SyncProducer

func Init(addr string) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Return.Errors = true
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Retry.Max = 5

	var err error
	producer, err = retryConnectionWithFixedDelay(saramaConfig, addr, 5, 2*time.Second)
	if err != nil {
		log.Fatalf("error creating kafka producer: %s", err)
		return
	}

	log.Printf("connected to Kafka on %s", addr)
}

func Send(topic string, values ...string) {
	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(strings.Join(values, " ")),
	}

	_, _, err := producer.SendMessage(message)
	if err != nil {
		log.Printf("Error sending message: %s", err)
	}
}

func retryConnectionWithFixedDelay(saramaConfig *sarama.Config, addr string, retries int, delay time.Duration) (sarama.SyncProducer, error) {
	for i := 0; i < retries; i++ {
		if producer, err := sarama.NewSyncProducer([]string{addr}, saramaConfig); err == nil {
			return producer, nil
		}
		time.Sleep(delay)
	}
	return nil, errors.New("failed to connect to Kafka")
}
