package main

import (
	"fmt"
	"log"
	"time"

	"context"

	"github.com/mohammadbnei/grpc-tuto/chat"
	"google.golang.org/grpc"
)

func main() {

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := chat.NewChatServiceClient(conn)

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("%dxHello From Client!", i)
		response, err := c.SayHello(context.Background(), &chat.Message{Body: message})
		if err != nil {
			log.Fatalf("Error when calling SayHello: %s", err)
		}
		log.Printf("Response from server: %s", response.Body)
		time.Sleep(time.Second * 3)
	}

}
