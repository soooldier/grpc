package main

import (
	"io"
	"log"
	"net"

	pb "github.com/soooldier/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type StreamService struct{}

func (s *StreamService) Server(r *pb.StreamRequest, stream pb.StreamService_ServerServer) error {
	for n := 0; n < 6; n++ {
		err := stream.Send(&pb.StreamResponse{
			Pt: &pb.StreamPoint{
				Name:  r.Pt.Name,
				Value: r.Pt.Value + int32(n),
			},
		})
		if err != nil {
			return nil
		}
	}
	return nil
}

func (s *StreamService) Client(stream pb.StreamService_ClientServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.StreamResponse{Pt: &pb.StreamPoint{
				Name:  "gRPC Stream Server: Client",
				Value: 1,
			}})
		}
		if err != nil {
			return err
		}
		log.Printf("stream.Recv => name: %s, value: %d", req.Pt.Name, req.Pt.Value)
	}
	return nil
}

func (s *StreamService) Both(stream pb.StreamService_BothServer) error {
	n := 0
	for {
		err := stream.Send(&pb.StreamResponse{
			Pt: &pb.StreamPoint{
				Name:  "gRPC stream Client: Both",
				Value: int32(n),
			},
		})
		if err != nil {
			return err
		}
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		log.Printf("Both recv => name: %s, value: %d", req.Pt.Name, req.Pt.Value)
		n++
	}
	return nil
}

const (
	PORT = "9002"
)

func main() {
	server := grpc.NewServer()
	pb.RegisterStreamServiceServer(server, &StreamService{})
	listen, _ := net.Listen("tcp", ":"+PORT)
	reflection.Register(server)
	server.Serve(listen)
}
