package worker

import (
	"context"
	"counter-service/internal/repository"
	"counter-service/pkg/kafka"
	"fmt"
	"log"
)

type WorkerI interface {
	LogRequestsEveryMinute()
}

type Worker struct {
	Repo repository.RepoI
}

func New(repo repository.RepoI) WorkerI {
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
