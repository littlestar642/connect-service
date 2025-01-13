package service

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func IsValidId(c *gin.Context) (int, error) {
	idStr := c.Query("id")
	if idStr == "" {
		return 0, errors.New("failed")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, errors.New("failed")
	}

	if id < 0 {
		return 0, errors.New("failed")
	}
	return id, nil
}
