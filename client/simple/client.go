package main

import (
	"context"
	pb "github.com/soooldier/grpc/proto"
	"google.golang.org/grpc"
	"log"
)

const PORT = "9001"

func main() {
	conn, err := grpc.Dial(":"+PORT, grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		log.Fatalf("grpc.Dail err: %v", err)
	}
	client := pb.NewSearchServiceClient(conn)
	resp, err := client.Search(context.Background(), &pb.SearchRequest{Request: "gRPC"})
	if err != nil {
		log.Fatalf("client.Search err: %v", err)
	}
	log.Printf("resp: %s", resp.GetResponse())
}
