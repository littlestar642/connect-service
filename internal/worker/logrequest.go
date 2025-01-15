package worker

import (
	"context"
	"counter-service/internal/repository"
	"counter-service/pkg/kafka"
	// "counter-service/pkg/logger"
	"fmt"
	"log"
)

type WorkerI interface {
	LogRequestsEveryMinute()
}

type Worker struct {
	Repo          repository.RepoI
	KafkaProducer kafka.ProducerI
}

func New(repo repository.RepoI, kafkaProducer kafka.ProducerI) WorkerI {
	return &Worker{
		Repo:          repo,
		KafkaProducer: kafkaProducer,
	}
}

func (w *Worker) LogRequestsEveryMinute() {
	log.Println("Logging request count")

	count, err := w.Repo.GetLastMinuteRequestCount(context.Background())
	if err != nil {
		log.Println("failed to get last minute request count: ", err)
		return
	}
	w.KafkaProducer.Send("request-count", fmt.Sprintf("%d", count))

	// logger.PrintToFile("Number of requests in the last minute: ", count)
}
