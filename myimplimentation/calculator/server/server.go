package main

import (
	"context"
	"grpc-go-course/myimplimentation/calculator/calculatorpb"
	"log"
	"net"

	"google.golang.org/grpc"
)

type CalulatorServer struct{}

func (*CalulatorServer) Sum(c context.Context, sr *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	log.Printf("Invoking sum service with value: %v", sr)
	res := &calculatorpb.SumResponse{
		Result: sr.GetA() + sr.GetB(),
	}

	return res, nil
}
func main() {
	log.Printf("Setting up a Server\n")
	cc, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failled to start Listining: %v\n", err)
	}
	s := grpc.NewServer()
	calculatorpb.RegisterCalulatorServer(s, &CalulatorServer{})
	if err := s.Serve(cc); err != nil {
		log.Fatalf("server failled with Error: %v\n", err)
	}

}
