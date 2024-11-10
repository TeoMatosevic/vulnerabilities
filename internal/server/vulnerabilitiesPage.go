package server

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) VulnerabilitiesPageHandler(c *gin.Context) {
	tmpl, err := template.ParseFiles("./internal/templates/vulnerabilities.gohtml")

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	cookie, err := generateNonce()

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("csrf-token", cookie, 3600, "/", "", false, true)

	tmpl.Execute(c.Writer, cookie)

	c.Writer.Header().Set("Content-Type", "text/html")
	c.Writer.WriteHeader(http.StatusOK)
}
