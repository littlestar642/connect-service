package service

import (
	"counter-service/internal/api"
	"counter-service/internal/repository"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type CounterI interface {
	Accept(c *gin.Context)
	AcceptCount(c *gin.Context)
}

type counterService struct {
	repo      repository.RepoI
	apiClient api.Requester
}

func New(repo repository.RepoI, apiClient api.Requester) CounterI {
	return &counterService{
		repo:      repo,
		apiClient: apiClient,
	}
}

func (s *counterService) Accept(c *gin.Context) {
	id, err := IsValidId(c)
	if err != nil {
		log.Printf("failed to validate id: %v\n", err)
		c.String(http.StatusBadRequest, "failed")
		return
	}

	if !s.repo.IsUniqueRequestId(c, id) {
		log.Printf("request id is not unique: %d\n", id)
		c.String(http.StatusConflict, "failed")
		return
	}

	err = s.repo.IncrementRequestCount(c, id)
	if err != nil {
		log.Printf("failed to increment request count: %v\n", err)
		c.String(http.StatusInternalServerError, "failed")
		return
	}

	endpoint := c.Query("endpoint")
	if endpoint != "" {
		decodedEndpoint, err := url.QueryUnescape(endpoint)
		if err != nil {
			log.Println("failed to decode endpoint:", err)
			c.String(http.StatusBadRequest, "failed")
			return
		}
		count, err := s.repo.GetCurrentMinuteRequestCount(c)
		if err != nil {
			log.Println("failed to get current minute request count:", err)
			c.String(http.StatusInternalServerError, "failed")
			return
		}

		go s.apiClient.SendPostRequest(decodedEndpoint, count)
	}

	c.String(http.StatusOK, "ok")
}

func (s *counterService) AcceptCount(c *gin.Context) {
	log.Println("recieved request for accept count")
}
