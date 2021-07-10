package main

import (
	"context"
	"grpc-go-course/myimplimentation/greet/greetpb"
	"io"
	"log"
	"sync"

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
	// doServerStream(c)
	// doClientStream(c)
	doBidiStream(c)
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

func doClientStream(c greetpb.GreetServiceClient) {
	log.Printf("Calling doClientStream")

	req := []*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Olivier",
				LastName:  "Barbier",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Junior",
				LastName:  "Barbier",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Stephanie",
				LastName:  "Barbier",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Deborah",
				LastName:  "Barbier",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while calling LongGreet: %v\n", err)
	}

	for _, r := range req {
		err := stream.Send(r)
		if err != nil {
			log.Fatalf("Error while sending Stream from LongGreet: %v\n", err)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while closing and Receiving Stream from LongGreet: %v\n", err)
	}

	log.Printf("Result for LongGreetRequest: %v\n", res.Result)

}

func doBidiStream(c greetpb.GreetServiceClient) {
	log.Printf("Calling doBidiStream")
	var wg sync.WaitGroup
	// var mu sync.RWMutex
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Failed to call GreetEveryone: %v\n", err)
	}
	req := []*greetpb.GreetEveryoneRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Olivier",
				LastName:  "Barbier",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Junior",
				LastName:  "Barbier",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Stephanie",
				LastName:  "Barbier",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Deborah",
				LastName:  "Barbier",
			},
		},
	}

	waitChan := make(chan struct{})
	sender := func(request []*greetpb.GreetEveryoneRequest) {
		for _, r := range request {
			wg.Add(1)
			go func(r *greetpb.GreetEveryoneRequest) {
				defer wg.Done()
				stream.Send(r)
			}(r)
		}

	}

	go func() {
		// waiting to close on different goroutine
		wg.Wait()
		stream.CloseSend()
	}()
	sender(req)

	receiver := func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Failed to received from streamer: %v\n", err)
				break
			}

			log.Printf("Result: %v\n", res.Result)
		}
		close(waitChan)
	}
	go receiver()
	<-waitChan
}
