package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Message struct {
	Type string  `json:"type"`
	IDs  []int64 `json:"ids"`
}

func PublishArray(ch *amqp.Channel, queue string, ids []int64, storyType string) {

	msg := Message{
		Type: storyType,
		IDs:  ids,
	}
	if len(ids) == 0 {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	body, _ := json.Marshal(msg)

	err := ch.PublishWithContext(ctx, "", queue, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	},
	)
	if err != nil {
		fmt.Println(" publish error:", err)
	} else {
		fmt.Printf("Published %d %s IDs to %s\n", len(ids), storyType, queue)
	}
}
