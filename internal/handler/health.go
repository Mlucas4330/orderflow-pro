package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthHandler struct {
	DB *pgxpool.Pool
}

func (h *HealthHandler) Check(c *gin.Context) {
	ctx := c.Request.Context()

	err := h.DB.Ping(ctx)
	if err != nil {
		log.Printf("Erro no health check do banco de dados: %v", err)

		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unhealthy",
			"error":  "database not available",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
