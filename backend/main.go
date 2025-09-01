package main

import (
	"backend/fetcher"
)

func main() {
	var oldTopIds, oldNewIds []int64
	ch := fetcher.ConnectToRabbitmq()
	fetcher.PeriodicFetcher(oldTopIds, oldNewIds, ch)
}
