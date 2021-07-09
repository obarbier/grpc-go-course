package main

import (
	"context"
	"grpc-go-course/myimplimentation/greet/greetpb"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {
	log.Println("Setting Client")

	cc, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer cc.Close()
	c := greetpb.NewGreetServiceClient(cc)
	// doUnary(c)
	doServerStream(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	log.Printf("Calling doUnary")
	greet := &greetpb.Greeting{
		FirstName: "Olivier",
		LastName:  "Barbier",
	}
	req := &greetpb.GreetRequest{
		Greeting: greet,
	}
	resp, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Response is not accepted: %v", err)
	}

	log.Printf("Response from Server: %v", resp.GetResult())
}

func doServerStream(c greetpb.GreetServiceClient) {
	log.Printf("Calling doServerStream")
	greet := &greetpb.Greeting{
		FirstName: "Olivier",
		LastName:  "Barbier",
	}
	req := &greetpb.GreetManyTimesRequest{
		Greeting: greet,
	}
	stream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error calling the GreetManyTimes RPC call")
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			log.Printf("End of the Stream, therefore Breaking")
			break
		}

		if err != nil {
			log.Fatalf("some error occured: %v", err)
		}

		log.Printf("Received Result from Stream Server: %v", res.Result)

	}

}
