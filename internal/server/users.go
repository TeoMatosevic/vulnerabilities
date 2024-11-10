package server

import (
	"github.com/gin-gonic/gin"
)

func (s *Server) GetUsersHandler(c *gin.Context) {
	firstName := c.Query("firstname")
	lastName := c.Query("lastname")

	query := "SELECT email FROM users WHERE first_name = $1 AND last_name = $2"

	rows, err := s.db.QueryResponable(query, firstName, lastName)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	jsonRows := make([]string, 0)

	for rows.Next() {
		var data string
		err := rows.Scan(&data)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		jsonRows = append(jsonRows, data)
	}

	c.JSON(200, jsonRows)
}
