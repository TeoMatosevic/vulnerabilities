package server

import (
	"github.com/gin-gonic/gin"
)

func (s *Server) GetUserHandler(c *gin.Context) {
	id, err := c.Cookie("session-id")

	if err != nil {
		c.JSON(400, gin.H{"error": "No session cookie found"})
		return
	}

	query := "SELECT first_name, last_name FROM users WHERE id = $1"

	rows, err := s.db.QueryResponable(query, id)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var firstName, lastName string

	for rows.Next() {
		err := rows.Scan(&firstName, &lastName)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(200, gin.H{"first_name": firstName, "last_name": lastName})
}
