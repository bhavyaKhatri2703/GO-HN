package main

import (
	"backend/application/api"
	g "backend/application/grpc_"
	"backend/fetcher"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func ConnectSQL() (*sql.DB, error) {
	connStr := "postgresql://postgres:postgres@localhost:5555/hackernews?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("db connected")
	return db, err
}

func main() {

	var oldTopIds []int64
	var oldNewIds []int64
	ch := fetcher.ConnectToRabbitmq()
	go fetcher.PeriodicFetcher(oldTopIds, oldNewIds, ch)

	db, err := ConnectSQL()
	if err != nil {
		fmt.Println("❌ Error connecting to the database:", err)
		return
	}

	grpcConn := g.Start_grpc()

	server := api.NewServer(db, grpcConn)
	server.SetRoutes()

	fmt.Println("Server is running on http://localhost:8080")
	if err := server.Router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
