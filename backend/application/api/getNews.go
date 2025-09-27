package api

import (
	search "backend/Search"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type GetNewsRequest struct {
	Names      []string  `json:"Names"`
	Embeddings []float32 `json:"Embeddings"`
}

func (s *Server) GetNews(ctx *gin.Context) {
	var req GetNewsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	bm25Query := strings.Join(req.Names, " ")
	embQuery := req.Embeddings

	topStories, newStories, err := search.HybridSearch(bm25Query, embQuery, s.DB)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving news"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":    "News retrieved successfully",
		"newStories": newStories,
		"topStories": topStories,
	})
}
