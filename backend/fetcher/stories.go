package fetcher

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type story struct {
	By    string `json:"by"`
	Id    int64  `json:"id"`
	Score int    `json:"score"`
	Time  int64  `json:"time"`
	Title string `json:"title"`
	Type  string `json:"type"`
	Url   string `json:"url"`
	Text  string `json:"text"`
}

func getStoriesIds(url string) []int64 {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("resp error")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var ids []int64
	json.Unmarshal(body, &ids)

	return ids
}

func getStory(id int64) story {
	url := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%s.json", strconv.FormatInt(id, 10))

	resp, _ := http.Get(url)

	var story story
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &story)
	return story

}

func toAddIds(old []int64, new []int64) []int64 {
	var add []int64
	mp := make(map[int64]int)
	for i := 0; i < len(old); i++ {
		mp[old[i]]++
	}
	for i := 0; i < len(new); i++ {
		if mp[new[i]] == 0 {
			add = append(add, new[i])
		}

	}

	return add
}

func toDeleteIds(old []int64, new []int64) []int64 {
	var delete []int64
	mp := make(map[int64]int)
	for i := 0; i < len(new); i++ {
		mp[new[i]]++
	}
	for i := 0; i < len(old); i++ {
		if mp[old[i]] == 0 {
			delete = append(delete, old[i])
		}
	}
	return delete
}

func ConnectToRabbitmq() *amqp.Channel {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("connceted to rmq")

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err, "hello")
	}

	queues := []string{"hn_add", "hn_delete"}
	for _, q := range queues {
		ch.QueueDeclare(q, true, false, false, false, nil)
	}

	return ch
}

func PeriodicFetcher(oldTopIds []int64, oldNewIds []int64, ch *amqp.Channel) {
	newStoriesUrl := "https://hacker-news.firebaseio.com/v0/newstories.json"
	topStoriesUrl := "https://hacker-news.firebaseio.com/v0/topstories.json"

	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		topIds := getStoriesIds(topStoriesUrl)
		newIds := getStoriesIds(newStoriesUrl)

		addTop := toAddIds(oldTopIds, topIds)
		deleteTop := toDeleteIds(oldTopIds, topIds)

		addNew := toAddIds(oldNewIds, newIds)
		deleteNew := toDeleteIds(oldNewIds, newIds)

		PublishArray(ch, "hn_add", addTop, "top")
		PublishArray(ch, "hn_delete", deleteTop, "top")
		PublishArray(ch, "hn_add", addNew, "new")
		PublishArray(ch, "hn_delete", deleteNew, "new")

		oldNewIds = newIds
		oldTopIds = topIds
	}

}
