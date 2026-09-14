package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var errBad = errors.New("invalid")

// badRequest: semua body-bind failure → 400 "invalid".
// Pesan validasi mentah (gin) jangan bocor ke klien (BUG-003).
func badRequest(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
}
