package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BasicAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, password, ok := c.Request.BasicAuth()
		if !ok || username != "adminbookly" || password != "adminbookly123" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Tidak diizinkan - Username atau Password salah",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}