package handler

import (
	"counter-service/internal/api"
	"counter-service/internal/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo *repository.Repo
}

func New(repo *repository.Repo) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) Accept(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		c.String(http.StatusBadRequest, "failed")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(http.StatusBadRequest, "failed")
		return
	}

	if id < 0 {
		c.String(http.StatusBadRequest, "failed")
		return
	}

	if !h.repo.IsUniqueId(c, id){
		c.String(http.StatusConflict, "failed")
		return
	}

	endpoint := c.Query("endpoint")
	if endpoint!=""{
		count := h.repo.GetCount(c)
		go api.SendPostRequest(endpoint, count)
	}

	c.String(http.StatusOK, "success")
}