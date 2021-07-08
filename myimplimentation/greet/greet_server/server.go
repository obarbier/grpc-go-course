package main

import (
	"fmt"
	"grpc-go-course/myimplimentation/greet/greetpb"
	"log"
	"net"

	"google.golang.org/grpc"
)

type Server struct{}

func main() {
	fmt.Println("Setting up GRPC Server")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, Server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to start Server: %v", err)
	}
}
