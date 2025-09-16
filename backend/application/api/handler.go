package api

import (
	"backend/application/grpc_"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Interests struct {
	Names []string `json:"names"`
}

type InterestsCookie struct {
	Names      []string  `json:"names"`
	Embeddings []float32 `json:"embeddings"`
}

func (s *Server) InterestsHandler(c *gin.Context) {
	var interests Interests

	if err := c.ShouldBindJSON(&interests); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid interests"})
		return
	}

	embeddings, err := grpc_.ReqEmbeddings(s.GrpcConn, interests.Names)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get embeddings"})
		return
	}

	cookieData := InterestsCookie{
		Names:      interests.Names,
		Embeddings: embeddings,
	}

	data, err := json.Marshal(cookieData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not encode data"})
		return
	}

	c.SetCookie("user_interests", string(data), 1800, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Interests and embeddings saved successfully"})
}
