package main

import (
	"fmt"
	"log"
	"time"

	"context"

	"github.com/mohammadbnei/grpc-tuto/person"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := person.NewPersonServiceClient(conn)

	for i := 0; i < 3; i++ {
		response, err := c.CreateUser(context.Background(), &person.Person{
			Name:        fmt.Sprintf("hugo %d", i),
			Email:       fmt.Sprintf("hugo%d@test.fr", i),
			LastUpdated: timestamppb.Now(),
		})
		if err != nil {
			log.Fatalf("Error when calling CreateUser: %s", err)
		}
		log.Printf("Response from server: %s", response.Id)
		time.Sleep(time.Second)
	}

	response, err := c.GetUsers(context.Background(), &person.Empty{})
	if err != nil {
		log.Fatalf("Error when calling GetUsers: %s", err)
	}
	log.Printf("Response from server: %s", response.People)

}
