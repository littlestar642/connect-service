package main

import (
	"context"
	"counter-service/internal/api"
	"counter-service/internal/config"
	"counter-service/internal/handler"
	"counter-service/internal/repository"
	"counter-service/internal/service"
	"counter-service/internal/worker"
	"counter-service/pkg/logger"
	"counter-service/pkg/redis"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	err := logger.Init()
	if err != nil {
		log.Fatalln("failed to open log file")
	}

	cnf := config.New()

	redisClient, err := redis.Init(cnf.RedisAddr)
	if err != nil {
		log.Fatalln("failed to connect to Redis: ", err.Error())
	}

	repo := repository.New(redisClient)

	wkr := worker.New(repo)

	apiClient := api.New()

	svc := service.New(repo, apiClient)

	handler := handler.New(svc)

	r := setupRouter(handler)

	taskTicker := time.NewTicker(time.Minute)
	done := make(chan bool)
	initTicker(wkr, taskTicker, done)

	r.Run(":" + cnf.Port)
	handleGracefulShutDown(&http.Server{Addr: ":" + cnf.Port, Handler: r}, 5*time.Second, taskTicker, done)
}

func setupRouter(handler *handler.Handler) *gin.Engine {
	r := gin.Default()
	r.GET("/api/verve/accept", handler.Accept)

	// adding post api for testing
	r.POST("/api/verve/accept", handler.AcceptCount)
	return r
}

func initTicker(wkr *worker.Worker, taskTicker *time.Ticker, done chan bool) {
	go func() {
		for {
			select {
			case <-done:
				log.Println("done from ticker")
				return
			case <-taskTicker.C:
				// handle error
				wkr.LogRequestsEveryMinute()
			}
		}
	}()
}

func handleGracefulShutDown(server *http.Server, timeout time.Duration, taskTicker *time.Ticker, done chan bool) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	done <- true
	taskTicker.Stop()
	_ = server.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}
