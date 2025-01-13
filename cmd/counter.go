package main

import (
	"counter-service/internal/config"
	"counter-service/internal/handler"
	"counter-service/internal/repository"
	"counter-service/internal/tasks"
	"counter-service/pkg/redis"
	"log"

	"github.com/gin-gonic/gin"
)

func main(){
	cnf := config.New()

	redisClient := redis.Init(cnf.RedisAddr)

	repo := repository.New(redisClient)

	handler := handler.New(repo)

	r:= gin.Default()
	r.GET("/api/verve/accept", handler.Accept)

	go tasks.LogRequestsEveryMinute()

	r.Run(":" + cnf.Port)
	log.Printf("Server started on port %s\n", cnf.Port)
}