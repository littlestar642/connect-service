package service

import (
	"counter-service/internal/api"
	"counter-service/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CounterService struct {
	repo *repository.Repo
}

func New(repo *repository.Repo) *CounterService {
	return &CounterService{
		repo: repo,
	}
}

func (s *CounterService) Accept(c *gin.Context) {
	id, err := IsValidId(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if !s.repo.IsUniqueRequestId(c, id) {
		c.String(http.StatusConflict, "failed")
		return
	}

	endpoint := c.Query("endpoint")
	if endpoint != "" {
		count := s.repo.GetRequestCount(c)
		go api.SendPostRequest(endpoint, count)
	}

	c.String(http.StatusOK, "ok")
}
