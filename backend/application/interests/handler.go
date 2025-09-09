package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Interests struct {
	Names      []string `json:"names"`
	Embeddings []int64  `json:"embeddings"`
}

func InterestsHandler(c *gin.Context, db *sql.DB) {
	var interests Interests

	err := c.ShouldBindJSON(&interests)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid interests"})
		return
	}
	interestsData, err := json.Marshal(interests)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not encode interests"})
		return
	}
	c.SetCookie("user_interests", string(interestsData), 1800, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Interests saved successfully"})
}
