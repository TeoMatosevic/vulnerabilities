package server

import (
	"crypto/sha256"
	"fmt"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Server) LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	query := fmt.Sprint("SELECT id, salt, pwd_hash FROM users WHERE email = $1")

	rows, err := s.db.QueryResponable(query, req.Email)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var id, salt, pwdHash string
	for rows.Next() {
		err := rows.Scan(&id, &salt, &pwdHash)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	if salt == "" {
		c.JSON(400, gin.H{"error": "User not found"})
		return
	}

	if pwdHash != hashPassword(req.Password, salt) {
		c.JSON(400, gin.H{"error": "Invalid password"})
		return
	}

	c.SetCookie("session-id", id, 3600, "/", "", false, true)

	c.JSON(200, gin.H{"message": "Login successful"})
}

func hashPassword(password, salt string) string {
	h := sha256.New()
	h.Write([]byte(salt + password))
	return fmt.Sprintf("%x", h.Sum(nil))
}
