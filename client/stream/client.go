package main

import (
	"context"
	"io"
	"log"

	pb "github.com/soooldier/grpc/proto"
	"google.golang.org/grpc"
)

const (
	PORT = "9002"
)

func main() {
	conn, _ := grpc.Dial(":"+PORT, grpc.WithInsecure())
	defer conn.Close()
	client := pb.NewStreamServiceClient(conn)
	printServer(client, &pb.StreamRequest{Pt: &pb.StreamPoint{Name: "gRPC Stream Client: Server", Value: 2018}})
	printClient(client, &pb.StreamRequest{Pt: &pb.StreamPoint{Name: "gRPC Stream Client: Client", Value: 2018}})
	printBoth(client, &pb.StreamRequest{Pt: &pb.StreamPoint{Name: "gRPC Stream Client: Both", Value: 2018}})
}

func printServer(client pb.StreamServiceClient, r *pb.StreamRequest) error {
	stream, err := client.Server(context.Background(), r)
	if err != nil {
		return err
	}
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		log.Printf("resp => name: %s, value: %d", resp.Pt.Name, resp.Pt.Value)
	}
	return nil
}

func printClient(client pb.StreamServiceClient, r *pb.StreamRequest) error {
	stream, err := client.Client(context.Background())
	if err != nil {
		return err
	}
	for n := 0; n < 6; n++ {
		if err := stream.Send(r); err != nil {
			return err
		}
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	log.Printf("resp => name: %s, value: %d", resp.Pt.Name, resp.Pt.Value)
	return nil
}

func printBoth(client pb.StreamServiceClient, r *pb.StreamRequest) error {
	stream, err := client.Both(context.Background())
	if err != nil {
		return err
	}
	for n := 0; n < 6; n++ {
		err := stream.Send(r)
		if err != nil {
			return err
		}
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		log.Printf("resp => name: %s, value: %d", resp.Pt.Name, resp.Pt.Value)
	}
	stream.CloseSend()
	return nil
}
