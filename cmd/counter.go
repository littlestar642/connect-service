package main

import (
	"counter-service/internal/config"
	"counter-service/internal/handler"
	"counter-service/internal/repository"
	"counter-service/internal/service"
	"counter-service/internal/worker"
	"counter-service/pkg/logger"
	"counter-service/pkg/redis"
	"log"

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
		logger.StdOut().Fatalln("Failed to connect to Redis: ", err.Error())
	}

	repo := repository.New(redisClient)

	wkr := worker.New(repo)

	svc := service.New(repo)

	handler := handler.New(svc)

	r := gin.Default()
	r.GET("/api/verve/accept", handler.Accept)

	go wkr.LogRequestsEveryMinute()

	r.Run(":" + cnf.Port)
	logger.StdOut().Println("Server started on port", cnf.Port)
}
