package main

import (
	"context"
	"fmt"
	pb "grpc-go-course/myimplimentation/blog/blogpb"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var dbConn *mongo.Client
var collection *mongo.Collection

type blogItem struct {
	// string id = 1;
	// string author_id = 2;
	// string title = 3;
	// string content = 4;
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Content  string             `bson:"content"`
	Title    string             `bson:"title"`
}

type blogServiceServer struct {
	// pb.UnimplementedBlogServiceServer
}

func (*blogServiceServer) ListBlog(req *pb.ListBlogRequest, stream pb.BlogService_ListBlogServer) error {
	fmt.Println("ListBlog request")
	cursor, err := collection.Find(context.Background(), bson.D{})
	defer cursor.Close(context.Background())
	if err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internall Error Occured: %v", err),
		)
	}
	for cursor.Next(context.Background()) {
		data := &blogItem{}
		if err := cursor.Decode(data); err != nil {
			return status.Errorf(
				codes.Internal,
				fmt.Sprintf("Cannot decode data: %v", err),
			)
		}
		stream.Send(&pb.ListBlogResponse{Blog: dataToBlogPb(data)})
	}
	if err := cursor.Err(); err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internall Error Occured: %v", err),
		)

	}
	return nil
}
func (*blogServiceServer) DeleteBlog(ctx context.Context, req *pb.DeleteBlogRequest) (*pb.DeleteBlogResponse, error) {
	fmt.Println("Delete blog request")
	oid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse ID"),
		)
	}

	filter := bson.M{"_id": oid}

	res, err := collection.DeleteOne(ctx, filter)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot delete object in MongoDB: %v", err),
		)
	}

	if res.DeletedCount == 0 {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find blog in MongoDB: %v", err),
		)
	}

	return &pb.DeleteBlogResponse{Id: req.GetId()}, nil
}

func (*blogServiceServer) UpdateBlog(ctx context.Context, req *pb.UpdateBlogRequest) (*pb.UpdateBlogResponse, error) {
	fmt.Println("Update blog request")
	blog := req.GetBlog()
	oid, err := primitive.ObjectIDFromHex(blog.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse ID"),
		)
	}

	// create an empty struct
	data := &blogItem{}
	filter := bson.M{"_id": oid}

	res := collection.FindOne(ctx, filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find blog with specified ID: %v", err),
		)
	}

	// we update our internal struct
	data.AuthorID = blog.GetAuthorId()
	data.Content = blog.GetContent()
	data.Title = blog.GetTitle()

	_, updateErr := collection.ReplaceOne(context.Background(), filter, data)
	if updateErr != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot update object in MongoDB: %v", updateErr),
		)
	}

	return &pb.UpdateBlogResponse{
		Blog: dataToBlogPb(data),
	}, nil
}

func (*blogServiceServer) ReadBlog(ctx context.Context, req *pb.ReadBlogRequest) (*pb.ReadBlogResponse, error) {
	log.Println("Invoked ReadBlog")
	id := req.GetId()
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("BlogID is not a valid ObjectID: %v\n", err)
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("BlogID is not a valid ObjectID: %v\n", err))
	}
	var data = &blogItem{}
	res := collection.FindOne(ctx, bson.M{"_id": oid})
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find blog with specified ID: %v", err),
		)
	}

	return &pb.ReadBlogResponse{
		Blog: dataToBlogPb(data),
	}, nil
}

func dataToBlogPb(data *blogItem) *pb.Blog {
	return &pb.Blog{
		Id:       data.ID.Hex(),
		AuthorId: data.AuthorID,
		Content:  data.Content,
		Title:    data.Title,
	}
}

func (*blogServiceServer) CreateBlog(ctx context.Context, req *pb.CreateBlogRequest) (*pb.CreateBlogResponse, error) {
	log.Printf("Create Blog Request")
	blog := req.GetBlog()
	document := &blogItem{
		AuthorID: blog.GetAuthorId(),
		Title:    blog.GetTitle(),
		Content:  blog.GetContent(),
	}

	res, err := collection.InsertOne(ctx, document)
	if err != nil {
		log.Printf("Internal Error: %v", err)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal Error: %v", err))
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Printf("Cannot Convert: %v", err)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Cannot Convert: %v", err))
	}

	return &pb.CreateBlogResponse{
		Blog: &pb.Blog{
			Id:       oid.Hex(),
			AuthorId: blog.GetAuthorId(),
			Title:    blog.GetTitle(),
			Content:  blog.GetContent(),
		},
	}, nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	cancel, dbConn, Dbctx := initDB()
	collection = dbConn.Database("mydb").Collection("blog")
	defer cancel()
	log.Printf("Setting up Blog Server")
	// NET
	l, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Printf("Cannot bing to provided network address: %v\n", err)
	}
	s := grpc.NewServer()
	// Bind service to grpc
	pb.RegisterBlogServiceServer(s, &blogServiceServer{})

	go func() {
		log.Printf("server listening at %v\n", l.Addr())
		if err := s.Serve(l); err != nil {
			log.Printf("Failed to serve: %v", err)
		}
	}()

	// Grcafully Shutting down server
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	fmt.Print("Disconnecting DB\n")
	dbConn.Disconnect(Dbctx)
	fmt.Printf("Stopping Blog Server\n")
	s.Stop()
	fmt.Printf("Clossing Listenner\n")
	l.Close()

}

func initDB() (context.CancelFunc, *mongo.Client, context.Context) {
	uri := "mongodb://localhost:27017"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected and pinged.")
	return cancel, client, ctx
}
