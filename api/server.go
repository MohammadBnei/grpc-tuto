package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/mohammadbnei/grpc-tuto/person"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://toor:passdorw@mongo/"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	defer client.Disconnect(ctx)

	db := client.Database("test")

	s := person.Server{
		Db: db,
	}

	grpcServer := grpc.NewServer()

	person.RegisterPersonServiceServer(grpcServer, &s)

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to listen on port 9000: %v", err)
	}

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server on port 9000: %v", err)
	}
}
