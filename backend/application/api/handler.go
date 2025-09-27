package api

import (
	"backend/application/grpc_"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Interests struct {
	Names []string `json:"names"`
}

type InterestsCookie struct {
	Names      []string  `json:"Names"`
	Embeddings []float32 `json:"Embeddings"`
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

	Data := InterestsCookie{
		Names:      interests.Names,
		Embeddings: embeddings,
	}

	c.JSON(http.StatusOK, Data)
}
