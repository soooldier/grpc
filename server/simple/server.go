package main

import (
	"context"
	"log"
	"net"

	pb "github.com/soooldier/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type SearchService struct{}

func (s *SearchService) Search(ctx context.Context, r *pb.SearchRequest) (*pb.SearchResponse, error) {
	return &pb.SearchResponse{Response: r.GetRequest() + " Server"}, nil
}

const PORT = "9001"

func main() {
	server := grpc.NewServer()
	pb.RegisterSearchServiceServer(server, &SearchService{})
	listen, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}
	reflection.Register(server)
	server.Serve(listen)
}
