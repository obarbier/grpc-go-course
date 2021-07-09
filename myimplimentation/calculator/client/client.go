package main

import (
	"context"
	"grpc-go-course/myimplimentation/calculator/calculatorpb"
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

	doSum(cClient)
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
