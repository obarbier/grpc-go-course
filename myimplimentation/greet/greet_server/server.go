package main

import (
	"context"
	"fmt"
	"grpc-go-course/myimplimentation/greet/greetpb"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	log.Printf("LongGreet function was invoked with Stream\n")
	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// Stream has closed therefore sending response
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})

		}

		if err != nil {
			log.Fatalf("Some error occureded in LongGreet: %v", err)
		}

		first_name := req.GetGreeting().GetFirstName()
		result += "Hello " + first_name + "! "
	}

}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Printf("Done Receing Data from client Stream\n")
			return nil
		}
		if err != nil {
			log.Fatalf("Error while receiving data from client stream: %v\n", err)
			return err
		}

		first_name := req.GetGreeting().GetFirstName()
		sendErr := stream.Send(&greetpb.GreetEveryoneResponse{
			Result: fmt.Sprintf("Hello %v !", first_name),
		})
		time.Sleep(2 * time.Second)
		if sendErr != nil {
			log.Fatalf("Error while Sending data to client: %v\n", sendErr)
			return sendErr
		}
	}
}
func (*server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineRespone, error) {
	log.Printf("GreetWithDeadline function was invoked with Request: %v\n", req)

	// for {
	select {
	case <-ctx.Done():
		//Timeout Excedded
		log.Print("The client canceled the request!")
		return nil, status.Error(codes.Canceled, "the client canceled the request")
	case <-time.After(3 * time.Second):
		firstname := req.GetGreeting().GetFirstName()
		response := fmt.Sprintf("Hello: %v", firstname)
		res := &greetpb.GreetWithDeadlineRespone{
			Result: response,
		}
		return res, nil
	}
	// }

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
