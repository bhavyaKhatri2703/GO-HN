package News

import (
	search "backend/Search"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type InterestsCookie struct {
	Names      []string  `json:"names"`
	Embeddings []float32 `json:"embeddings"`
}

func GetNews(ctx *gin.Context, db *sql.DB) {
	cookie, err := ctx.Cookie("user_interests")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No user interests found"})
		return
	}

	var cookieData InterestsCookie
	err = json.Unmarshal([]byte(cookie), &cookieData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid cookie data"})
		return
	}

	bm25Query := strings.Join(cookieData.Names, " ")
	embQuery := cookieData.Embeddings

	stories, err := search.HybridSearch(bm25Query, embQuery, db)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving news"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "News retrieved successfully",
		"stories": stories,
	})
}
