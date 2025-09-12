package search

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

type Story struct {
	Id        int64  `json:"id"`
	By        string `json:"by"`
	Type      string `json:"type"`
	Text      string `json:"text"`
	Url       string `json:"url"`
	Title     string `json:"title"`
	Full_text string `json:"full_text"`
	Score     int    `json:"score"`
}

func HybridSearch(query string, queryEmb []float32, db *sql.DB) ([]Story, error) {
	var stories []Story

	// Convert []float32 to a properly formatted string for the PostgreSQL vector type
	embStr := make([]string, len(queryEmb))
	for i, v := range queryEmb {
		embStr[i] = fmt.Sprintf("%f", v)
	}
	embeddings := fmt.Sprintf("[%s]", strings.Join(embStr, ","))

	rows, err := db.Query(`
    SELECT id, by, type, text, url, title, full_text, score,
           bm25_embedding <&> to_bm25query('top_embedding_bm25', tokenize($1 , 'tokenizer1')) AS bm25_rank,
           (sem_embedding <=> $2::vector) AS distance
    FROM topStories
    ORDER BY ((bm25_embedding <&> to_bm25query('top_embedding_bm25', tokenize($1 , 'tokenizer1'))) * 0.5) +
             ((sem_embedding <=> $2::vector) * 0.5)
    LIMIT 20
`, query, embeddings)

	if err != nil {
		log.Println("Query error:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var story Story
		var bm25_rank, distance float32
		err := rows.Scan(&story.Id, &story.By, &story.Type, &story.Text, &story.Url, &story.Title, &story.Full_text, &story.Score, &bm25_rank, &distance)
		if err != nil {
			log.Println("Scan error:", err)
			return nil, err
		}
		fmt.Println(story)
		stories = append(stories, story)
	}
	return stories, nil
}
