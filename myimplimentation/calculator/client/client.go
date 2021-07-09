package main

import (
	"context"
	"grpc-go-course/myimplimentation/calculator/calculatorpb"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Fail to Dial: %v", err)
	}

	cClient := calculatorpb.NewCalulatorClient(conn)
	defer conn.Close()

	// doSum(cClient)
	doPrimeNumberDecomposition(cClient)
}

func doSum(cClient calculatorpb.CalulatorClient) {
	req := &calculatorpb.SumRequest{
		A: 5,
		B: 10,
	}
	res, err := cClient.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to calculate Sum: %v\n", err)
	}
	log.Printf("Sum is: %v\n", res.Result)
}

func doPrimeNumberDecomposition(cClient calculatorpb.CalulatorClient) {
	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 120,
	}

	stream, err := cClient.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to invoke PrimeNumberDecomposition: %v ", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			log.Printf("Stream ended therefore breaking")
			break
		}

		if err != nil {
			log.Fatalf("Some Error occured: %v", err)
		}

		log.Printf("Result from PrimeNumberDecomposition is: %v ", res.Result)
	}

}
