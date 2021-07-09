package main

import (
	"context"
	"fmt"
	"grpc-go-course/myimplimentation/greet/greetpb"
	"log"
	"net"

	"google.golang.org/grpc"
)

type Server struct{}

func (*Server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	log.Printf("greet function was invoked with Request: %v\n", req)
	firstname := req.GetGreeting().GetFirstName()
	response := fmt.Sprintf("Hello: %v", firstname)
	res := &greetpb.GreetResponse{
		Result: response,
	}
	return res, nil
}

func main() {
	fmt.Println("Setting up GRPC Server")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &Server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to start Server: %v", err)
	}
}
