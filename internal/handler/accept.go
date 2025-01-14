package handler

import (
	"counter-service/internal/service"

	"github.com/gin-gonic/gin"
)

type HandlerI interface{
	Accept(c *gin.Context)
	AcceptCount(c *gin.Context)
}
type handler struct {
	svc service.CounterI
}

func New(serv service.CounterI) HandlerI {
	return &handler{
		svc: serv,
	}
}

func (h *handler) Accept(c *gin.Context) {
	h.svc.Accept(c)
}

func (h *handler) AcceptCount(c *gin.Context) {
	h.svc.AcceptCount(c)
}
