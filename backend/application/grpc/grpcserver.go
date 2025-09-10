package grpc

import (
	pb "backend/proto"
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
)

func Start_grpc() *grpc.ClientConn {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	return conn
}

func ReqEmbeddings(conn *grpc.ClientConn) *pb.InterestsResponse {
	client := pb.NewEmbeddingsServiceClient(conn)
	req := &pb.InterestsRequest{
		Interests: []string{"python", "Go", "Ai"}, // will come from cookies
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.GetEmbeddings(ctx, req)
	if err != nil {
		log.Fatalf("Error calling service: %v", err)
	}
	return res
}
