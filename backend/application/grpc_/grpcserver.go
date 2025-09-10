package grpc_

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

func ReqEmbeddings(conn *grpc.ClientConn, interests []string) ([]float32, error) {
	client := pb.NewEmbeddingsServiceClient(conn)
	req := &pb.InterestsRequest{
		Interests: interests,
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.GetEmbeddings(ctx, req)
	if err != nil {
		log.Fatalf("Error calling service: %v", err)
	}
	return res.Embeddings, err
}
