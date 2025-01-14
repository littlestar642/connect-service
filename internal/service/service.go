package service

import (
	"counter-service/internal/api"
	"counter-service/internal/repository"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CounterService struct {
	repo      *repository.Repo
	apiClient *api.API
}

func New(repo *repository.Repo, apiClient *api.API) *CounterService {
	return &CounterService{
		repo:      repo,
		apiClient: apiClient,
	}
}

func (s *CounterService) Accept(c *gin.Context) {
	id, err := IsValidId(c)
	if err != nil {
		log.Printf("Failed to validate id: %v\n", err)
		c.String(http.StatusBadRequest, "failed")
		return
	}

	if !s.repo.IsUniqueRequestId(c, id) {
		log.Printf("Request id is not unique: %d\n", id)
		c.String(http.StatusConflict, "failed")
		return
	}

	err = s.repo.IncrementRequestCount(c, id)
	if err != nil {
		log.Printf("Failed to increment request count: %v\n", err)
		c.String(http.StatusInternalServerError, "failed")
		return
	}

	endpoint := c.Query("endpoint")
	if endpoint != "" {
		count := s.repo.GetRequestCount(c)
		go s.apiClient.SendPostRequest(endpoint, count)
	}

	c.String(http.StatusOK, "ok")
}
