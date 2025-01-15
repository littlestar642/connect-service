package main

import (
	"context"
	"counter-service/internal/api"
	"counter-service/internal/config"
	"counter-service/internal/handler"
	"counter-service/internal/repository"
	"counter-service/internal/service"
	"counter-service/internal/worker"
	"counter-service/pkg/kafka"
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
		log.Fatalln("failed to open log file", err.Error())
	}

	cnf := config.New()

	redisClient, err := redis.Init(cnf.RedisAddr)
	if err != nil {
		log.Fatalln("failed to connect to redis: ", err.Error())
	}

	kafkaClient, err := kafka.InitProducer(cnf.KafkaAddr)
	if err != nil {
		log.Fatalln("failed to connect to kafka: ", err.Error())
	}

	repo := repository.New(redisClient)

	wkr := worker.New(repo, kafkaClient)

	apiClient := api.New()

	svc := service.New(repo, apiClient)

	handler := handler.New(svc)

	r := setupRouter(handler)

	taskTicker := time.NewTicker(time.Minute)
	done := make(chan bool)
	initTicker(wkr, taskTicker, done)

	err = r.Run(":" + cnf.Port)
	if err != nil {
		log.Fatalln("failed to start server: ", err.Error())
	}

	log.Println("server started")
	handleGracefulShutDown(&http.Server{Addr: ":" + cnf.Port, Handler: r}, 5*time.Second, taskTicker, done, kafkaClient, redisClient)
}

func setupRouter(handler handler.HandlerI) *gin.Engine {
	r := gin.Default()
	r.GET("/api/verve/accept", handler.Accept)

	// adding post api for testing
	r.POST("/api/verve/accept", handler.AcceptCount)
	return r
}

func initTicker(wkr worker.WorkerI, taskTicker *time.Ticker, done chan bool) {
	go func() {
		for {
			select {
			case <-done:
				log.Println("done from ticker")
				return
			case <-taskTicker.C:
				wkr.LogRequestsEveryMinute()
			}
		}
	}()
}

func handleGracefulShutDown(server *http.Server, timeout time.Duration, taskTicker *time.Ticker, done chan bool, kafkaClient kafka.ProducerI, redisClient redis.ClientI) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	done <- true
	taskTicker.Stop()
	kafkaClient.Close()
	redisClient.Close()
	_ = server.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}
