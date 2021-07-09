package main

import (
	"context"
	"grpc-go-course/myimplimentation/greet/greetpb"
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
	doUnary(c)
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
