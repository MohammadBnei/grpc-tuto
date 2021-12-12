# gRPC Tutorial - Person CRUD

## Protobuf

Create a person folder, inside it the person.proto
```
    syntax = "proto3";
    package person;

    option go_package = "/person";

    message Person {
        string name = 1;
        string id = 2;  // Unique ID number for this person.
        string email = 3;

        enum PhoneType {
            MOBILE = 0;
            HOME = 1;
            WORK = 2;
        }

        message PhoneNumber {
            string number = 1;
            PhoneType type = 2;
        }

        repeated PhoneNumber phones = 4;
    }

    // Our address book file is just one of these.
    message AddressBook {
        repeated Person people = 1;
    }

    message Id {
        string id = 1;
    }

    message Empty {}

    service PersonService {
        rpc createUser (Person) returns (Id) {}
        rpc getUsers (Empty) returns (AddressBook) {}
    }
```

Then you generate the interfaces
```
    protoc --go_out=plugins=grpc:. person/person.proto
```

## CRUD with gRPC

In the server.go file, at these lines at the top of the main function
```
    client, err := mongo.NewClient(options.Client().ApplyURI("<MONGO_URL>"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	defer client.Disconnect(ctx)

	db := client.Database("<MONGO_DATABASE>")
```

Add the database instance in the gRPC server struct  
```
    type Server struct {
        Db *mongo.Database
    }
```

And create the appropriate functions
```
    func (s *Server) CreateUser(ctx context.Context, person *Person) (*Id, error) {
        inserted, err := s.Db.Collection("person").InsertOne(ctx, person)
        if err != nil {
            log.Fatalf("Error inserting person: %v", err)
            return nil, err
        }
        id := &Id{
            Id: fmt.Sprint(inserted.InsertedID),
        }
        return id, nil
    }

    func (s *Server) GetUsers(ctx context.Context, _ *Empty) (*AddressBook, error) {
        cursor, err := s.Db.Collection("person").Find(ctx, bson.M{})
        if err != nil {
            log.Fatalf("Error finding all persons: %v", err)
            return nil, err
        }

        addressBook := &AddressBook{}
        if err := cursor.All(ctx, &addressBook.People); err != nil {
            log.Fatalf("Error decoding persons: %v", err)
            return nil, err
        }

        return addressBook, nil
    }
```

Finally, update the client to test the functions
```
    c := person.NewPersonServiceClient(conn)

	for i := 0; i < 3; i++ {
		response, err := c.CreateUser(context.Background(), &person.Person{
			Name:  fmt.Sprintf("hugo %d", i),
			Email: fmt.Sprintf("hugo%d@test.fr", i),
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
```

And voila !
```
    go run client/client.go
```