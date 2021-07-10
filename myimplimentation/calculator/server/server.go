package main

import (
	"context"
	"fmt"
	"grpc-go-course/myimplimentation/calculator/calculatorpb"
	"io"
	"log"
	"math"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CalulatorServer struct{}

func (*CalulatorServer) Sum(c context.Context, sr *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	log.Printf("Invoking sum service with value: %v\n", sr)
	res := &calculatorpb.SumResponse{
		Result: sr.GetA() + sr.GetB(),
	}

	return res, nil
}

func (*CalulatorServer) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.Calulator_PrimeNumberDecompositionServer) error {
	log.Printf("Invoking PrimeNumberDecomposition service with value: %v\n", req)
	N := req.GetNumber()
	var K int32 = 2
	for N > 1 {
		if N%K == 0 {
			stream.Send(&calculatorpb.PrimeNumberDecompositionResponse{
				Result: K,
			})
			N = N / K
		} else {
			K++
		}
	}
	return nil
}

func (*CalulatorServer) ComputeAverage(stream calculatorpb.Calulator_ComputeAverageServer) error {
	log.Printf("Invoking ComputeAverage service with stream of number\n")
	count := 0
	var sum float32 = 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// End of the stream therefore needs to return
			res := &calculatorpb.ComputeAverageResponse{
				Result: sum / float32(count),
			}
			return stream.SendAndClose(res)
		}

		if err != nil {
			log.Fatalf("Some error occured while recieving from client: %v\n", err)
		}
		sum += req.GetNumber()
		count++
	}
}

func (*CalulatorServer) FindMaximum(stream calculatorpb.Calulator_FindMaximumServer) error {
	log.Printf("Invoking FindMaximum service with stream of number\n")
	var max float64
	numprocess := 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Printf("Done processing the stream")
			return nil
		}
		if err != nil {
			log.Fatalf("Failed to receive from client stream: %v\n", err)
			return err
		}
		if numprocess == 0 {
			numprocess++
			max = req.GetNumber()
			sendErr := stream.Send(&calculatorpb.FindMaximumResponse{
				Result: max,
			})
			if sendErr != nil {
				log.Fatalf("Error while sending to the client stream: %v\n", sendErr)
				return sendErr
			}
		} else {
			numprocess++
			if max < req.GetNumber() {
				max = req.GetNumber()
				sendErr := stream.Send(&calculatorpb.FindMaximumResponse{
					Result: max,
				})
				if sendErr != nil {
					log.Fatalf("Error while sending to the client stream: %v\n", sendErr)
					return sendErr
				}
			}
		}
	}

}

func (*CalulatorServer) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	log.Printf("Invoking sum service with value: %v\n", req)
	number := req.GetNumber()

	if number < 0 {
		return nil, status.Errorf(codes.InvalidArgument,
			fmt.Sprintf("Receiver a negative number: %d\n", number))
	}

	return &calculatorpb.SquareRootResponse{
		Result: math.Sqrt(float64(number)),
	}, nil
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
