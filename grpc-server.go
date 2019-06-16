package main

import (
	context "context"
	"net"

	pb "github.com/smirnov/grpc-echo/pb"
	"google.golang.org/grpc"
)

type EchoServer struct {
}

//Simple echoing functionality returning inbound message as is
func (s EchoServer) Echo(ctx context.Context, inbound *pb.Message) (*pb.Message, error) {
	return inbound, nil
}

func main() {
	echoServer := EchoServer{}
	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterEchoServiceServer(grpcServer, echoServer)
	grpcServer.Serve(listen)
}
