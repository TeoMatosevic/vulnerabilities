package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func (s *Server) VulnerableUsersHandler(c *gin.Context) {
	firstName := c.Query("firstname")
	lastName := c.Query("lastname")

	query := fmt.Sprintf("SELECT email FROM users WHERE first_name = '%s' AND last_name = '%s'", firstName, lastName)

	rows, err := s.db.Query(query)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	jsonRows := make([]*interface{}, 0)
	for rows.Next() {
		var data *interface{}
		err := rows.Scan(&data)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		jsonRows = append(jsonRows, data)
	}

	c.JSON(200, jsonRows)
}
