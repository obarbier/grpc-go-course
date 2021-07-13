package main

import (
	"context"
	"fmt"
	pb "grpc-go-course/myimplimentation/blog/blogpb"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {
	log.Printf("Setting up Blog Server Client")

	conn, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewBlogServiceClient(conn)
	log.Printf("Creating Blog")
	blog := &pb.Blog{
		AuthorId: "obarbier",
		Title:    "my first Blog",
		Content:  "Hello world",
	}
	res, err := client.CreateBlog(context.Background(), &pb.CreateBlogRequest{Blog: blog})
	blogID := res.GetBlog().GetId()
	if err != nil {
		log.Printf("Unexpected Error: %v\n", err)
	}
	log.Printf("Created Blog: %v", res)

	//ReadBlog
	_, err2 := client.ReadBlog(context.Background(), &pb.ReadBlogRequest{
		Id: "fkjlkjfsd",
	})
	if err2 != nil {
		log.Printf("%v", err2)
	}

	res2, err2 := client.ReadBlog(context.Background(), &pb.ReadBlogRequest{
		Id: blogID,
	})
	if err2 != nil {
		log.Printf("%v", err2)
	}
	fmt.Println(res2)

	// update Blog
	newBlog := &pb.Blog{
		Id:       blogID,
		AuthorId: "Changed Author",
		Title:    "My First Blog (edited)",
		Content:  "Content of the first blog, with some awesome additions!",
	}
	updateRes, updateErr := client.UpdateBlog(context.Background(), &pb.UpdateBlogRequest{Blog: newBlog})
	if updateErr != nil {
		fmt.Printf("Error happened while updating: %v \n", updateErr)
	}
	fmt.Printf("Blog was updated: %v\n", updateRes)

	// delete Blog
	deleteRes, deleteErr := client.DeleteBlog(context.Background(), &pb.DeleteBlogRequest{Id: blogID})

	if deleteErr != nil {
		fmt.Printf("Error happened while deleting: %v \n", deleteErr)
	}
	fmt.Printf("Blog was deleted: %v \n", deleteRes)

	// list Blogs

	stream, err := client.ListBlog(context.Background(), &pb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("error while calling ListBlog RPC: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happened: %v", err)
		}
		fmt.Println(res.GetBlog())
	}

}
