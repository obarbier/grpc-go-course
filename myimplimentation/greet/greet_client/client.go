package main

import (
	"fmt"
	"grpc-go-course/myimplimentation/greet/greetpb"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Setting Client")

	cc, err := grpc.Dial("localhost:500501", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer cc.Close()
	c := greetpb.NewGreetServiceClient(cc)

	fmt.Printf("created client:%v", c)
}
