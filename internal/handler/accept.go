package handler

import (
	"counter-service/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *service.CounterService
}

func New(serv *service.CounterService) *Handler {
	return &Handler{
		svc: serv,
	}
}

func (h *Handler) Accept(c *gin.Context) {
	h.svc.Accept(c)
}

func (h *Handler) AcceptCount(c *gin.Context) {
	h.svc.AcceptCount(c)
}
