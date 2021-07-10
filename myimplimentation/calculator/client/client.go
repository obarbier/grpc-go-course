package main

import (
	"context"
	"grpc-go-course/myimplimentation/calculator/calculatorpb"
	"io"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	conn, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Fail to Dial: %v", err)
	}

	cClient := calculatorpb.NewCalulatorClient(conn)
	defer conn.Close()

	// doSum(cClient)
	// doPrimeNumberDecomposition(cClient)
	// doComputeAverage(cClient)
	// doBidi(cClient)
	doErrorCode(cClient)
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

func doComputeAverage(cClient calculatorpb.CalulatorClient) {
	req := []*calculatorpb.ComputeAverageRequest{
		{
			Number: 1,
		},
		{
			Number: 2,
		},
		{
			Number: 3,
		},
		{
			Number: 4,
		},
	}
	stream, err := cClient.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Failed to connect to ComputeAverage: %v\n", err)
	}
	for _, r := range req {
		err := stream.Send(r)
		if err != nil {
			log.Fatalf("Failed to Send stream to backend: %v\n", err)
		}

	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving resullt from backend: %v\n", err)
	}

	log.Printf("ComputeAverage Result: %v\n", res.Result)

}

func doBidi(cClient calculatorpb.CalulatorClient) {
	reqs := []float64{-1, -5, -3, -6, -2, -20}

	stream, err := cClient.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Failed to call FindMaximum: %v\n", err)
	}

	for _, r := range reqs {
		stream.Send(&calculatorpb.FindMaximumRequest{
			Number: r,
		})
	}

	stream.CloseSend()

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			log.Printf("Done Proccessing data from server stream")
			break
		}
		if err != nil {
			log.Fatalf("Error occured while receiving data from server stream: %v\n", err)
			return
		}
		log.Printf("Result: %v\n", res.Result)
	}

}

func doErrorCode(cClient calculatorpb.CalulatorClient) {

	call := func(number int32) {
		res, err := cClient.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{
			Number: number,
		})
		if err != nil {
			respErr, ok := status.FromError(err)
			if ok {
				// error comming from GRPC
				log.Printf("Grpc Error: %v\n", respErr.Message())
				if respErr.Code() == codes.InvalidArgument {
					log.Printf("Negative number was sent")
				}
			} else {
				// Framework error
				log.Fatalf("Fatal error:%v\n", respErr)
			}

			return
		}

		log.Printf("Result: %v\n", res.Result)
	}

	call(4)
	call(-4)

}
