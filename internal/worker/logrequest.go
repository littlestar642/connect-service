package worker

import (
	"context"
	"counter-service/internal/repository"
	"counter-service/pkg/logger"
	"time"
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
	time.Sleep(time.Minute)

	count := w.Repo.GetRequestCount(context.Background())
	logger.PrintToFile("Number of requests in the last minute: ", count)
}
