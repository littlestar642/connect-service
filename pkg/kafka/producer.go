package kafka

import (
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/IBM/sarama"
)

type ProducerI interface {
	Send(topic string, values ...string)
	Close()
}

type message struct {
	Topic     string
	Value     string
	Timestamp string
}

type producer struct {
	client sarama.SyncProducer
}

func InitProducer(addr string) (ProducerI, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Return.Errors = true
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Retry.Max = 5

	var err error
	producer, err := retryConnectionWithFixedDelay(saramaConfig, addr, 5, 2*time.Second)
	if err != nil {
		return nil, err
	}

	log.Printf("connected to Kafka on %s", addr)
	return producer, nil
}

func (p *producer) Send(topic string, values ...string) {
	msg := &message{
		Topic:     topic,
		Value:     strings.Join(values, " "), // delimeter can be chosen based on the requirement
		Timestamp: time.Now().String(),
	}

	jsonMsg, _ := json.Marshal(msg)
	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(jsonMsg),
	}

	_, _, err := p.client.SendMessage(message)
	if err != nil {
		log.Printf("Error sending message: %s", err)
	}
}

func retryConnectionWithFixedDelay(saramaConfig *sarama.Config, addr string, retries int, delay time.Duration) (ProducerI, error) {
	for i := 0; i < retries; i++ {
		if client, err := sarama.NewSyncProducer([]string{addr}, saramaConfig); err == nil {
			return &producer{
				client: client,
			}, nil
		}
		time.Sleep(delay)
	}
	return nil, errors.New("failed to connect to Kafka")
}

func (p *producer) Close() {
	if p.client != nil {
		p.client.Close()
	}
}
