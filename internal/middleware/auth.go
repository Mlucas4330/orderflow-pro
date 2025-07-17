package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		mockCustomerID, _ := uuid.Parse("00000000-0000-0000-0000-000000000001")

		c.Set("userID", mockCustomerID)

		c.Next()
	}
}
