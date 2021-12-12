# gRPC Tutorial

## Go

Launch the docker-compose file and open a terminal inside the app service
```
    docker-compose up -d
    docker-compose exec app bash
```

Then you will init a module using the path to your repository
```
    go mod init github.com/$USERNAME/$REPOSITORY
```

 You can then create a server.go file and copy the next lines
```
    package main

    import (
        "fmt"
        "time"
    )

    func main() {
        for {
            fmt.Println("Hello World")
            time.Sleep(time.Second * 3)
        }
    }
```

Initialize the live-reloader [air](https://github.com/cosmtrek/air) and start the app
```
    air init
    air

```
It should be displaying Hello Wolrd every 3 seconds. Try changing the string and see the live reload.


## Protobuf

Create a chat folder, and inside it a chat.proto file
We will define our service for gRPC
```
    syntax = "proto3";
    package chat;

    option go_package = "/";

    message Message {
        string body = 1;
    }

    // The greeting service definition.
    service ChatService {
        // Sends a greeting
        rpc SayHello (Message) returns (Message) {}
    }
```

To generate very useful interfaces and struct from the proto file, we can use the protoc cli command  
```
    # Install the protobuf packages
    go get -u github.com/golang/protobuf/protoc-gen-go

    # Install the protobuf compiler
    apt update && apt install protobuf-compiler

    # Generate the go files
    protoc --go_out=plugins=grpc:chat chat/chat.proto
```

Add the protobuf server to the app, change the main.go file
```
    import (
        "log"
        "net"

        "google.golang.org/grpc"
    )

    func main() {
        lis, err := net.Listen("tcp", ":9000")
        if err != nil {
            log.Fatalf("Failed to listen on port 9000: %v", err)
        }

        grpcServer := grpc.NewServer()
        if err := grpcServer.Serve(lis); err != nil {
            log.cupFatalf("Failed to serve gRPC server on port 9000: %v", err)
        }
    }
```

Go into you host terminal and restart the container, it will automatically start air
```
    docker-compose up --force-recreate app
```

### Simple gRPC implementation

Create a chat/chat.go file, and copy the next lines
```
    package chat

    import (
        "context"
        "log"
    )

    type Server struct {
    }

    func (s *Server) SayHello(ctx context.Context, in *Message) (*Message, error) {
        log.Printf("Receive message body from client: %s", in.Body)
        return &Message{Body: "Hello From the server"}, nil
    }

```

Add the chat gRPC to the server 
```
    s := chat.Server{}

    #Existing command
	grpcServer := grpc.NewServer()

	chat.RegisterChatServiceServer(grpcServer, &s)
```

Create the client server, client/client.go
```
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
```

You can start the client by typing the following (inside the container):
```
    go run client/client.go
```

It should be displaying " xHello From Client!", and "Hello From the server!"