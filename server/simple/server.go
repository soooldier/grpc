package main

import (
	"context"
	"log"
	"net/http"
	"runtime/debug"
	"strings"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	pb "github.com/soooldier/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SearchService struct{}

func (s *SearchService) Search(ctx context.Context, r *pb.SearchRequest) (*pb.SearchResponse, error) {
	return &pb.SearchResponse{Response: r.GetRequest() + " Server"}, nil
}

const PORT = "9001"

func RecoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			debug.PrintStack()
			err = status.Errorf(codes.Internal, "Panic error: %v", e)
		}
	}()
	return handler(ctx, req)
}

func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	log.Printf("gRPC method: %s, %v", info.FullMethod, req)
	resp, err = handler(ctx, req)
	log.Printf("gRPC method: %s, %v", info.FullMethod, resp)
	return resp, err
}

func main() {
	opts := []grpc.ServerOption{
		grpc_middleware.WithUnaryServerChain(
			RecoveryInterceptor,
			LoggingInterceptor,
		),
	}
	server := grpc.NewServer(opts...)
	pb.RegisterSearchServiceServer(server, &SearchService{})
	http.ListenAndServe(":"+PORT, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-type"), "application/grpc") {
			server.ServeHTTP(w, r)
		} else {
		}
	}))


	//listen, err := net.Listen("tcp", ":"+PORT)
	//if err != nil {
	//	log.Fatalf("net.Listen err: %v", err)
	//}
	//reflection.Register(server)
	//server.Serve(listen)
}
