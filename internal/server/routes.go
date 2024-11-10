package server

import (
	"crypto/rand"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.GET("/", s.VulnerabilitiesPageHandler)
	r.GET("/api/v1/users", s.GetUsersHandler)
	r.GET("/api/v1/vulnerable-users", s.VulnerableUsersHandler)
	r.POST("/api/v1/login", s.LoginHandler)
	r.GET("/api/v1/user", csrfMiddleware(), s.GetUserHandler)

	return r
}

func generateNonce() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}

func csrfMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		v := c.Query("vulnerable")

		if v == "true" {
			c.Next()
			return
		}

		token, err := c.Cookie("csrf-token")
		if err != nil {
			c.JSON(403, gin.H{"error": "No CSRF token found"})
			c.Abort()
		}

		if token != c.GetHeader("X-CSRF-Token") {
			c.JSON(403, gin.H{"error": "CSRF token mismatch"})
			c.Abort()
		}

		c.Next()
	}
}
