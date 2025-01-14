package main

import (
	"context"
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
		log.Fatalln("Failed to open log file")
	}

	cnf := config.New()

	redisClient, err := redis.Init(cnf.RedisAddr)
	if err != nil {
		log.Fatalln("Failed to connect to Redis: ", err.Error())
	}

	repo := repository.New(redisClient)

	wkr := worker.New(repo)

	svc := service.New(repo)

	handler := handler.New(svc)

	r := gin.Default()
	r.GET("/api/verve/accept", handler.Accept)

	done := make(chan bool)
	initTicker(wkr, done)

	r.Run(":" + cnf.Port)
	handleGracefulShutDown(&http.Server{Addr: ":" + cnf.Port, Handler: r}, 5*time.Second, done)
}

func initTicker(wkr *worker.Worker, done chan bool) {
	taskTicker := time.NewTicker(time.Minute)
	go func() {
		select {
		case <-done:
			log.Println("done from ticker")
			return
		case <-taskTicker.C:
			// handle error
			wkr.LogRequestsEveryMinute()
		}
	}()
}

func handleGracefulShutDown(server *http.Server, timeout time.Duration, done chan bool) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	done <- true
	_ = server.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}
