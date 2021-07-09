package main

import (
	"context"
	"fmt"
	"grpc-go-course/myimplimentation/greet/greetpb"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	log.Printf("greet function was invoked with Request: %v\n", req)
	firstname := req.GetGreeting().GetFirstName()
	response := fmt.Sprintf("Hello: %v", firstname)
	res := &greetpb.GreetResponse{
		Result: response,
	}
	return res, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	log.Printf("GreetManyTimes function was invoked with Request: %v\n", req)
	for i := 0; i < 10; i++ {
		res := &greetpb.GreetManyTimesResponse{
			Result: fmt.Sprintf("Hello %v for %d times \n", req.GetGreeting().GetFirstName(), i),
		}

		stream.Send(res)
		time.Sleep(2 * time.Second)

	}

	return nil
}

func main() {
	fmt.Println("Setting up GRPC Server")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to start Server: %v", err)
	}
}
