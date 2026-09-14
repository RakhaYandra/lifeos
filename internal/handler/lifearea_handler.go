package handler

import (
	"net/http"

	"github.com/RakhaYandra/lifeos/internal/repository"
	"github.com/gin-gonic/gin"
)

type LifeAreaHandler struct {
	Areas *repository.LifeAreaRepository
}

func (h *LifeAreaHandler) List(c *gin.Context) {
	out, err := h.Areas.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	if out == nil {
		out = []repository.LifeAreaRow{}
	}
	c.JSON(http.StatusOK, out)
}
