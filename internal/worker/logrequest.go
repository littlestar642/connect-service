package worker

import (
	"context"
	"counter-service/internal/repository"
	"counter-service/pkg/kafka"
	"fmt"
	"log"
)

type Worker struct {
	Repo *repository.Repo
}

func New(repo *repository.Repo) *Worker {
	return &Worker{
		Repo: repo,
	}
}

func (w *Worker) LogRequestsEveryMinute() {
	log.Println("Logging request count")

	count := w.Repo.GetLastMinuteRequestCount(context.Background())
	kafka.Send("request-count", fmt.Sprintf("%d", count))
	// logger.PrintToFile("Number of requests in the last minute: ", count)
}
